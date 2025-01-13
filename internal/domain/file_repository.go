package domain

import (
	"context"
	"fmt"

	"github.com/daniarmas/notes/internal/clog"
	"github.com/daniarmas/notes/internal/oss"
	"github.com/google/uuid"
)

type FileRepository interface {
	Create(ctx context.Context, ossFileId, path string, noteID uuid.UUID) (*File, error)
	Update() error
	Delete() error
	List() error
	Move() error
	Process(ctx context.Context, ossFileId string) error
}

type fileCloudRepository struct {
	FileDatabaseDs       FileDatabaseDs
	ObjectStorageService oss.ObjectStorageService
}

func NewFileRepository(fileDatabaseDs FileDatabaseDs, objectStorageService oss.ObjectStorageService) FileRepository {
	return &fileCloudRepository{
		FileDatabaseDs:       fileDatabaseDs,
		ObjectStorageService: objectStorageService,
	}
}

func (r *fileCloudRepository) Create(ctx context.Context, ossFileId, path string, noteID uuid.UUID) (*File, error) {
	if ossFileId != "" {
		// Save the file on the database
		file := &File{OriginalFile: ossFileId, NoteId: noteID}
		file, err := r.FileDatabaseDs.CreateFile(ctx, file)
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

func (r *fileCloudRepository) Delete() error { return nil }

func (r *fileCloudRepository) List() error { return nil }

func (r *fileCloudRepository) Move() error { return nil }

func (r *fileCloudRepository) Process(ctx context.Context, ossFileId string) error {
	// Download the file from the cloud
	path, err := r.ObjectStorageService.GetObject(ctx, ossFileId)
	if err != nil {
		return err
	}
	clog.Info(ctx, fmt.Sprintf("path: %v", path), nil)
	return nil
}
