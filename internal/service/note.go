package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/daniarmas/clogg"
	"github.com/daniarmas/notes/internal/config"
	"github.com/daniarmas/notes/internal/customerrors"
	"github.com/daniarmas/notes/internal/domain"
	"github.com/daniarmas/notes/internal/k8sc"
	"github.com/daniarmas/notes/internal/oss"
	"github.com/google/uuid"
)

type CreateNoteResponse struct {
	Note *domain.Note `json:"note"`
}

type PresignedUrl struct {
	Url      string `json:"url"`
	File     string `json:"file"`
	ObjectId string `json:"object_id"`
}

// GetPresignedUrlsResponse represents the structure of the get presigned urls response
type GetPresignedUrlsResponse struct {
	Urls []PresignedUrl `json:"urls"`
}

type NoteService interface {
	CreateNote(ctx context.Context, title string, content string, objectNames []string) (*CreateNoteResponse, error)
	ListTrashNotesByUser(ctx context.Context, cursor time.Time) (*[]domain.Note, error)
	ListNotesByUser(ctx context.Context, cursor time.Time) (*[]domain.Note, error)
	RestoreNote(ctx context.Context, id uuid.UUID) (*domain.Note, error)
	DeleteNote(ctx context.Context, id uuid.UUID, hard bool) error
	UpdateNote(ctx context.Context, note *domain.Note) (*domain.Note, error)
	GetPresignedUrls(ctx context.Context, objectNames []string) (*GetPresignedUrlsResponse, error)
}

type noteService struct {
	Config         config.Configuration
	FileRepository domain.FileRepository
	NoteRepository domain.NoteRepository
	Oss            oss.ObjectStorageService
	K8sClient      k8sc.K8sC
	Db             *sql.DB
}

func NewNoteService(noteRepository domain.NoteRepository, oss oss.ObjectStorageService, fileRepository domain.FileRepository, cfg config.Configuration, k8sClient k8sc.K8sC, db *sql.DB) NoteService {
	return &noteService{
		NoteRepository: noteRepository,
		Oss:            oss,
		FileRepository: fileRepository,
		Config:         cfg,
		K8sClient:      k8sClient,
		Db:             db,
	}
}

func (s *noteService) CreateNote(ctx context.Context, title string, content string, objectNames []string) (*CreateNoteResponse, error) {
	// Start the sql transaction
	tx, err := s.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Defer the transaction rollback or commit
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	note := &domain.Note{
		UserId:  domain.GetUserIdFromContext(ctx),
		Title:   title,
		Content: content,
	}

	// Check concurrently if the objects exists in the oss
	var wg sync.WaitGroup
	errChan := make(chan error, len(objectNames))
	for _, objectName := range objectNames {
		wg.Add(1)
		go func(objectName string) {
			defer wg.Done()
			err := s.Oss.ObjectExists(ctx, s.Config.ObjectStorageServiceBucket, objectName)
			if err != nil {
				errChan <- err
				return
			}
		}(objectName)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, errors.New("objects not found")
	}

	// Create the note
	note, err = s.NoteRepository.CreateNote(ctx, tx, note)
	if err != nil {
		return nil, err
	}
	// Create the files concurrently
	var files []*domain.File
	var mu2 sync.Mutex
	var wg2 sync.WaitGroup
	errChan2 := make(chan error, len(objectNames))
	for _, objectName := range objectNames {
		wg2.Add(1)
		go func(objectName string) {
			defer wg2.Done()
			file, err := s.FileRepository.Create(ctx, tx, objectName, "", note.Id)
			if err != nil {
				errChan2 <- err
				return
			}
			mu2.Lock()
			files = append(files, file)
			mu2.Unlock()
		}(objectName)
	}

	wg2.Wait()
	close(errChan2)

	if len(errChan2) > 0 {
		return nil, errors.New("error creating files")
	}

	// Create a k8s job to process the files
	if s.Config.InK8s {
		jobName := fmt.Sprintf("process-note-files-job-%s", note.Id)
		namespace := "default"
		imageName := s.Config.DockerImageName
		args := []string{
			"process-files",
			"--files",
		}
		// Append the slice of object names as a comma-separated string
		args = append(args, strings.Join(objectNames, ","))

		// Define the environment variables for the job
		envs := []corev1.EnvFromSource{
			{
				SecretRef: &corev1.SecretEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "note-secrets",
					},
				},
			},
		}

		err := s.K8sClient.CreateJob(ctx, jobName, namespace, imageName, args, envs)
		if err != nil {
			clogg.Error(ctx, "error creating k8s job", clogg.String("error", err.Error()))
			return nil, err
		}
	} else {
		// This is a mock for the k8s job on dev environment
		for _, file := range objectNames {
			if err := s.FileRepository.Process(ctx, tx, file); err != nil {
				clogg.Error(ctx, "error processing file", clogg.String("error", err.Error()))
			}
		}
	}

	// Generate the presigned urls to get the original files concurrently
	var mu3 sync.Mutex
	var wg3 sync.WaitGroup
	errChan3 := make(chan error, len(files))
	for _, file := range files {
		wg3.Add(1)
		go func(file *domain.File) {
			defer wg3.Done()
			url, err := s.Oss.PresignedGetObject(ctx, s.Config.ObjectStorageServiceBucket, file.OriginalFile, time.Second*24*60*60)
			if err != nil {
				errChan3 <- err
				return
			}
			mu3.Lock()
			file.Url = url
			mu3.Unlock()
		}(file)
	}

	wg3.Wait()
	close(errChan3)

	if len(errChan3) > 0 {
		return nil, errors.New("error getting the presigned urls to get the original files")
	}

	// Include the files in the note
	note.Files = files

	return &CreateNoteResponse{Note: note}, nil
}

func (s *noteService) ListNotesByUser(ctx context.Context, cursor time.Time) (*[]domain.Note, error) {
	// Get the user ID from the context
	userId := domain.GetUserIdFromContext(ctx)

	// Get the notes
	notes, err := s.NoteRepository.ListNotesByUser(ctx, userId, cursor)
	if err != nil {
		return nil, err
	}

	// Get all the ids from notes
	ids := make([]uuid.UUID, len(*notes))
	for i, note := range *notes {
		ids[i] = note.Id
	}

	// Get the files for each note
	files, err := s.FileRepository.ListFilesByNotesIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	// Generate the presigned urls to get the files
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, len(*files))
	for i := range *files {
		wg.Add(1)
		go func(file *domain.File) {
			defer wg.Done()
			var objectName string
			if file.ProcessedFile != "" {
				objectName = file.ProcessedFile
			} else {
				objectName = file.OriginalFile
			}
			url, err := s.Oss.PresignedGetObject(ctx, s.Config.ObjectStorageServiceBucket, objectName, time.Second*24*60*60)
			if err != nil {
				errChan <- err
				return
			}
			mu.Lock()
			file.Url = url
			mu.Unlock()
		}(&(*files)[i])
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, errors.New("error getting the presigned urls")
	}

	// Include the files in the notes
	// Create a map to group files by NoteId
	fileMap := make(map[uuid.UUID][]*domain.File)
	for _, file := range *files {
		fileMap[file.NoteId] = append(fileMap[file.NoteId], &file)
	}

	// Include the files in the notes
	for i, note := range *notes {
		if noteFiles, ok := fileMap[note.Id]; ok {
			(*notes)[i].Files = noteFiles
		}
	}

	return notes, nil
}

func (s *noteService) ListTrashNotesByUser(ctx context.Context, cursor time.Time) (*[]domain.Note, error) {
	// Get the user ID from the context
	userId := domain.GetUserIdFromContext(ctx)

	// Get the notes
	notes, err := s.NoteRepository.ListTrashNotesByUser(ctx, userId, cursor)
	if err != nil {
		return nil, err
	}

	// Get all the ids from notes
	ids := make([]uuid.UUID, len(*notes))
	for i, note := range *notes {
		ids[i] = note.Id
	}

	// Get the files for each note
	files, err := s.FileRepository.ListFilesByNotesIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	// Generate the presigned urls to get the files
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, len(*files))
	for i := range *files {
		wg.Add(1)
		go func(file *domain.File) {
			defer wg.Done()
			var objectName string
			if file.ProcessedFile != "" {
				objectName = file.ProcessedFile
			} else {
				objectName = file.OriginalFile
			}
			url, err := s.Oss.PresignedGetObject(ctx, s.Config.ObjectStorageServiceBucket, objectName, time.Second*24*60*60)
			if err != nil {
				errChan <- err
				return
			}
			mu.Lock()
			file.Url = url
			mu.Unlock()
		}(&(*files)[i])
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, errors.New("error getting the presigned urls")
	}

	// Include the files in the notes
	// Create a map to group files by NoteId
	fileMap := make(map[uuid.UUID][]*domain.File)
	for _, file := range *files {
		fileMap[file.NoteId] = append(fileMap[file.NoteId], &file)
	}

	// Include the files in the notes
	for i, note := range *notes {
		if noteFiles, ok := fileMap[note.Id]; ok {
			(*notes)[i].Files = noteFiles
		}
	}

	return notes, nil
}

func (s *noteService) RestoreNote(ctx context.Context, id uuid.UUID) (*domain.Note, error) {
	// Start the sql transaction
	tx, err := s.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Defer the transaction rollback or commit
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	note, err := s.NoteRepository.RestoreNote(ctx, tx, id)
	if err != nil {
		switch err.(type) {
		case *customerrors.RecordNotFound:
			return nil, errors.New("note not found")
		}
	}

	return note, nil
}

func (s *noteService) UpdateNote(ctx context.Context, note *domain.Note) (*domain.Note, error) {
	// Start the sql transaction
	tx, err := s.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Defer the transaction rollback or commit
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	note, err = s.NoteRepository.UpdateNote(ctx, tx, note)
	if err != nil {
		switch err.(type) {
		case *customerrors.RecordNotFound:
			return nil, errors.New("note not found")
		}
	}

	return note, nil
}

func (s *noteService) DeleteNote(ctx context.Context, id uuid.UUID, isHard bool) error {
	// Start the sql transaction
	tx, err := s.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Defer the transaction rollback or commit
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	var files *[]domain.File

	// Get the files if isHard is true before they are deleted from the database
	if isHard {
		if files, err = s.FileRepository.ListFilesByNoteId(ctx, id); err != nil {
			return err
		}
	}

	if err = s.NoteRepository.DeleteNote(ctx, tx, id, isHard); err != nil {
		if _, ok := err.(*customerrors.RecordNotFound); ok {
			return errors.New("note not found")
		}
		return err
	}

	// Delete the files from the cloud
	if isHard {
		if err = s.FileRepository.HardDeleteFiles(ctx, tx, files); err != nil {
			return err
		}
	}

	return nil
}

func (s *noteService) GetPresignedUrls(ctx context.Context, objectNames []string) (*GetPresignedUrlsResponse, error) {
	// Make a slice of presigned urls
	urls := make([]PresignedUrl, 0, len(objectNames))

	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, len(objectNames))

	for _, objectName := range objectNames {
		wg.Add(1)
		go func(objectName string) {
			defer wg.Done()
			// Generate a new object name
			id := uuid.New()
			ext := filepath.Ext(objectName)
			newObjectName := fmt.Sprintf("original/%s%s", id, ext)
			// Generate the presigned url
			url, err := s.Oss.PresignedPutObject(ctx, s.Config.ObjectStorageServiceBucket, newObjectName)
			if err != nil {
				errChan <- err
				return
			}
			mu.Lock()
			urls = append(urls, PresignedUrl{Url: url, File: objectName, ObjectId: newObjectName})
			mu.Unlock()
		}(objectName)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	return &GetPresignedUrlsResponse{Urls: urls}, nil
}
