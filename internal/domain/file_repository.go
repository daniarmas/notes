package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"image/jpeg"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/daniarmas/notes/internal/clog"
	"github.com/daniarmas/notes/internal/config"
	"github.com/daniarmas/notes/internal/oss"
	"github.com/google/uuid"
)

var pictureExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	// ".png":  true,
	// ".gif":  true,
	// ".bmp":  true,
	// ".tiff": true,
	// ".svg":  true,
}

var audioExtensions = map[string]bool{
	// ".mp3":  true,
	// ".wav":  true,
	// ".flac": true,
	// ".aac":  true,
	// ".ogg":  true,
	".m4a": true,
	// ".wma":  true,
}

// Returns string value if the file is a picture or audio
func fileType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if pictureExtensions[ext] {
		return "picture"
	} else if audioExtensions[ext] {
		return "audio"
	} else {
		return "unsupported"
	}
}

type FileRepository interface {
	Create(ctx context.Context, tx *sql.Tx, ossFileId, path string, noteID uuid.UUID) (*File, error)
	Update() error
	HardDeleteFiles(ctx context.Context, tx *sql.Tx, files *[]File) error
	ListFilesByNoteId(ctx context.Context, noteId uuid.UUID) (*[]File, error)
	ListFilesByNotesIds(ctx context.Context, noteId []uuid.UUID) (*[]File, error)
	Move() error
	Process(ctx context.Context, tx *sql.Tx, ossFileId string) error
}

type fileCloudRepository struct {
	FileDatabaseDs       FileDatabaseDs
	ObjectStorageService oss.ObjectStorageService
	Config               *config.Configuration
}

func NewFileRepository(fileDatabaseDs FileDatabaseDs, objectStorageService oss.ObjectStorageService, cfg *config.Configuration) FileRepository {
	return &fileCloudRepository{
		FileDatabaseDs:       fileDatabaseDs,
		ObjectStorageService: objectStorageService,
		Config:               cfg,
	}
}

func (r *fileCloudRepository) Create(ctx context.Context, tx *sql.Tx, ossFileId, path string, noteID uuid.UUID) (*File, error) {
	if ossFileId != "" {
		// Save the file on the database
		file := &File{OriginalFile: ossFileId, NoteId: noteID}
		file, err := r.FileDatabaseDs.CreateFile(ctx, tx, file)
		if err != nil {
			return nil, err
		}
		return file, nil
	} else if path != "" {
		// Save the file on the database and the cloud
		return nil, nil
	} else {
		// Return error
		return nil, nil
	}
}

func (r *fileCloudRepository) Update() error { return nil }

func (r *fileCloudRepository) HardDeleteFiles(ctx context.Context, tx *sql.Tx, files *[]File) error {
	processedFilesNames := make([]string, 0, len(*files))
	originalFilesNames := make([]string, 0, len(*files))
	// Get the processed and original files names
	for _, file := range *files {
		processedFilesNames = append(processedFilesNames, file.ProcessedFile)
		originalFilesNames = append(originalFilesNames, file.OriginalFile)
	}
	// Delete the files from the cloud
	var wg sync.WaitGroup
	errChan := make(chan error, len(processedFilesNames))
	for i, _ := range processedFilesNames {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processedFile := processedFilesNames[i]
			originalFile := originalFilesNames[i]
			err := r.ObjectStorageService.RemoveObject(ctx, r.Config.ObjectStorageServiceBucket, processedFile)
			if err != nil {
				errChan <- errors.New(fmt.Sprintf("error removing processed file %s from the cloud", processedFile))
				return
			}
			err = r.ObjectStorageService.RemoveObject(ctx, r.Config.ObjectStorageServiceBucket, originalFile)
			if err != nil {
				errChan <- errors.New(fmt.Sprintf("error removing original file %s from the cloud", originalFile))
				return
			}
		}()
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return errors.New("error removing files from the cloud")
	}

	return nil
}

func (r *fileCloudRepository) ListFilesByNotesIds(ctx context.Context, noteId []uuid.UUID) (*[]File, error) {
	// Fetch the files from the database
	files, err := r.FileDatabaseDs.ListFilesByNotesIds(ctx, noteId)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (r *fileCloudRepository) ListFilesByNoteId(ctx context.Context, noteId uuid.UUID) (*[]File, error) {
	// Fetch the file ids from the database
	files, err := r.FileDatabaseDs.ListFilesByNoteId(ctx, noteId)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (r *fileCloudRepository) Move() error { return nil }

func (r *fileCloudRepository) Process(ctx context.Context, tx *sql.Tx, ossFileId string) error {
	// Declare the processed file id
	var processedFileId string

	// Download the file from the cloud
	path, err := r.ObjectStorageService.GetObject(ctx, r.Config.ObjectStorageServiceBucket, ossFileId)
	if err != nil {
		return err
	}

	// Remove the file from tmp after the process
	defer os.Remove(path)

	// Process the file based on the type
	switch fileType(path) {
	case "picture":
		if path, err = CompressJpegImage(path); err != nil {
			return err
		}
		processedFileId = fmt.Sprintf("processed-photos/%s", filepath.Base(path))
	case "audio":
		if path, err = CompressAudio(path); err != nil {
			return err
		}
		processedFileId = fmt.Sprintf("processed-audio/%s", filepath.Base(path))
	default:
		clog.Error(ctx, "The file is not supported", nil)
		return fmt.Errorf("the file is not supported")
	}

	// Remove the processed file from tmp after the upload
	defer os.Remove(path)

	// Upload the processed file to the cloud
	if err := r.ObjectStorageService.PutObject(ctx, r.Config.ObjectStorageServiceBucket, processedFileId, path); err != nil {
		clog.Error(ctx, "error uploading processed file to the cloud", err)
		return err
	}

	// Update the file on the database
	if _, err := r.FileDatabaseDs.UpdateFileByOriginalId(ctx, tx, ossFileId, processedFileId); err != nil {
		clog.Error(ctx, "error updating file by original id", err)
		return err
	}

	return nil
}

// CompressAudio compresses the audio file using ffmpeg
func CompressAudio(path string) (string, error) {
	// Define the output path for the compressed audio
	id := uuid.New()
	ext := filepath.Ext(path)
	outputPath := fmt.Sprintf("/tmp/%s%s", id, ext)

	// Compress the audio using ffmpeg
	cmd := exec.Command("ffmpeg", "-i", path, "-b:a", "128k", outputPath)

	// Run the command
	if err := cmd.Run(); err != nil {
		clog.Error(context.Background(), "error compressing audio", err)
		return "", fmt.Errorf("error compressing audio: %v", err)
	}

	return outputPath, nil
}

// CompressJPEGImage compresses the image file using "image/jpeg"
func CompressJpegImage(path string) (string, error) {
	// Define the output path for the compressed audio
	id := uuid.New()
	ext := filepath.Ext(path)
	outputPath := fmt.Sprintf("/tmp/%s%s", id, ext)

	// Open the input file
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open input file: %v", err)
	}
	defer file.Close()

	// Decode the image
	img, err := jpeg.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode JPEG image: %v", err)
	}

	// Create the output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	// Encode the image with the specified quality
	options := &jpeg.Options{Quality: 50}
	if err := jpeg.Encode(outFile, img, options); err != nil {
		// Remove the output file if the encoding fails
		os.Remove(outputPath)

		return "", fmt.Errorf("failed to encode JPEG image: %v", err)
	}

	return outputPath, nil
}
