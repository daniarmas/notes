package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/daniarmas/notes/internal/clog"
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
}

func NewNoteService(noteRepository domain.NoteRepository, oss oss.ObjectStorageService, fileRepository domain.FileRepository, cfg config.Configuration, k8sClient k8sc.K8sC) NoteService {
	return &noteService{
		NoteRepository: noteRepository,
		Oss:            oss,
		FileRepository: fileRepository,
		Config:         cfg,
		K8sClient:      k8sClient,
	}
}

func (s *noteService) CreateNote(ctx context.Context, title string, content string, objectNames []string) (*CreateNoteResponse, error) {
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
	note, err := s.NoteRepository.CreateNote(ctx, note)
	if err != nil {
		return nil, err
	}
	// Create the files concurrently
	var wg2 sync.WaitGroup
	errChan2 := make(chan error, len(objectNames))
	for _, objectName := range objectNames {
		wg2.Add(1)
		go func(objectName string) {
			defer wg2.Done()
			_, err := s.FileRepository.Create(ctx, objectName, "", note.Id)
			if err != nil {
				errChan2 <- err
				return
			}
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

		err := s.K8sClient.CreateJob(ctx, jobName, namespace, imageName, args)
		if err != nil {
			clog.Error(ctx, "error creating k8s job", err)
			return nil, err
		}
	} else {
		// This is a mock for the k8s job on dev environment
		for _, file := range objectNames {
			if err := s.FileRepository.Process(ctx, file); err != nil {
				clog.Error(ctx, "error processing file", err)
			}
		}
	}

	return &CreateNoteResponse{Note: note}, nil
}

func (s *noteService) ListNotesByUser(ctx context.Context, cursor time.Time) (*[]domain.Note, error) {
	// Get the user ID from the context
	userId := domain.GetUserIdFromContext(ctx)

	notes, err := s.NoteRepository.ListNotesByUser(ctx, userId, cursor)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (s *noteService) ListTrashNotesByUser(ctx context.Context, cursor time.Time) (*[]domain.Note, error) {
	// Get the user ID from the context
	userId := domain.GetUserIdFromContext(ctx)

	notes, err := s.NoteRepository.ListTrashNotesByUser(ctx, userId, cursor)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (s *noteService) RestoreNote(ctx context.Context, id uuid.UUID) (*domain.Note, error) {
	note, err := s.NoteRepository.RestoreNote(ctx, id)
	if err != nil {
		switch err.(type) {
		case *customerrors.RecordNotFound:
			return nil, errors.New("note not found")
		}
	}
	return note, nil
}

func (s *noteService) UpdateNote(ctx context.Context, note *domain.Note) (*domain.Note, error) {
	note, err := s.NoteRepository.UpdateNote(ctx, note)
	if err != nil {
		switch err.(type) {
		case *customerrors.RecordNotFound:
			return nil, errors.New("note not found")
		}
	}
	return note, nil
}

func (s *noteService) DeleteNote(ctx context.Context, id uuid.UUID, isHard bool) error {
	err := s.NoteRepository.DeleteNote(ctx, id, isHard)
	if err != nil {
		switch err.(type) {
		case *customerrors.RecordNotFound:
			return errors.New("note not found")
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
			url, err := s.Oss.GetPresignedUrl(ctx, s.Config.ObjectStorageServiceBucket, newObjectName)
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
