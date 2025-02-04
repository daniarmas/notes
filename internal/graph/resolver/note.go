package resolver

import (
	"context"
	"errors"
	"time"

	"github.com/daniarmas/notes/internal/domain"
	"github.com/daniarmas/notes/internal/graph/model"
	"github.com/daniarmas/notes/internal/service"
	"github.com/daniarmas/notes/internal/utils"
	"github.com/google/uuid"
)

func mapFile(file domain.File) *model.File {
	var updateTime string
	if !file.UpdateTime.IsZero() {
		updateTime = file.UpdateTime.Format(time.RFC3339)

	}
	return &model.File{
		ID:            file.Id.String(),
		NoteID:        file.NoteId.String(),
		OriginalFile:  file.OriginalFile,
		ProcessedFile: &file.ProcessedFile,
		URL:           file.Url,
		CreateTime:    file.CreateTime.Format(time.RFC3339),
		UpdateTime:    &updateTime,
	}
}

// func to map domain.Note to model.Note
func mapNote(note domain.Note) *model.Note {
	var updateTime string
	if !note.UpdateTime.IsZero() {
		updateTime = note.UpdateTime.Format(time.RFC3339)
	}
	// map domain.File to model.File
	files := make([]*model.File, len(note.Files))
	for i, file := range note.Files {
		files[i] = mapFile(*file)
	}
	return &model.Note{
		ID:         note.Id.String(),
		UserID:     note.UserId.String(),
		Title:      &note.Title,
		Content:    &note.Content,
		Files:      files,
		CreateTime: note.CreateTime.Format(time.RFC3339),
		UpdateTime: &updateTime,
	}
}

// ListNotes is the resolver for the notes field.
func ListNotes(ctx context.Context, input *model.NotesInput, srv service.NoteService) (*model.NotesResponse, error) {
	// Check if the user is authenticated
	userId := domain.GetUserIdFromContext(ctx)
	if userId == uuid.Nil {
		return nil, errors.New("unauthenticated")
	}

	// Get the cursor from the query parameters
	var cursorQueryParam string
	if input != nil && input.Cursor != nil {
		cursorQueryParam = *input.Cursor
	}
	// parse the cursor query parameter
	if cursorQueryParam == "" {
		cursorQueryParam = time.Now().UTC().Format(time.RFC3339)
	}
	cursor, err := utils.ParseTime(cursorQueryParam)
	if err != nil && cursorQueryParam != "" {
		msg := "Invalid time format for the cursor query parameter. Must use RFC3339 format"
		return nil, errors.New(msg)
	}

	// Check if trash is true
	var notes *[]domain.Note

	if input != nil && input.Trash != nil && *input.Trash {
		notes, err = srv.ListTrashNotesByUser(ctx, cursor)
	} else {
		notes, err = srv.ListNotesByUser(ctx, cursor)
	}

	if err != nil {
		switch err.Error() {
		default:
			return nil, errors.New("internal server error")
		}
	}

	// Get the next cursor
	notesSlice := *notes
	var nextCursor time.Time
	if len(notesSlice) > 0 {
		nextCursor = notesSlice[len(notesSlice)-1].UpdateTime
	} else {
		// Handle the case where notesSlice is empty
		nextCursor = time.Now().UTC()
	}

	// Parse []domain.Note to []*model.Note
	notesRes := make([]*model.Note, len(notesSlice))
	for i, note := range notesSlice {
		notesRes[i] = mapNote(note)
	}

	return &model.NotesResponse{
		Notes:  notesRes,
		Cursor: nextCursor.Format(time.RFC3339),
	}, nil
}

// CreateNote is the resolver for the createNote field.
func CreateNote(ctx context.Context, input model.CreateNoteInput, srv service.NoteService) (*model.Note, error) {
	// Check if the user is authenticated
	userId := domain.GetUserIdFromContext(ctx)
	if userId == uuid.Nil {
		return nil, errors.New("unauthenticated")
	}

	var (
		title   string
		content string
	)

	if input.Title != nil {
		title = *input.Title
	}

	if input.Content != nil {
		content = *input.Content
	}

	objectNames := make([]string, len(input.ObjectNames))
	for i, objectName := range input.ObjectNames {
		objectNames[i] = *objectName
	}

	// Validate the input
	if input.Title == nil || *input.Title == "" {
		return nil, errors.New("field 'title' is required")
	}

	res, err := srv.CreateNote(ctx, title, content, objectNames)
	if err != nil {
		switch err.Error() {
		case "objects not found":
			msg := "One or more objects not found in the object storage service"
			return nil, errors.New(msg)
		default:
			return nil, errors.New("internal server error")
		}
	}
	// Parse domain.Note to model.Note
	return mapNote(*res.Note), nil
}

// CreatePresignedURL is the resolver for the createPresignedUrl field.
func CreatePresignedURL(ctx context.Context, objectName []string, srv service.NoteService) (*model.CreatePresignedUrlsResponse, error) {
	// Check if the user is authenticated
	userId := domain.GetUserIdFromContext(ctx)
	if userId == uuid.Nil {
		return nil, errors.New("unauthenticated")
	}

	// Validate the input
	if len(objectName) == 0 {
		return nil, errors.New("field 'objectName' is required")
	}

	res, err := srv.GetPresignedUrls(ctx, objectName)
	if err != nil {
		return nil, errors.New("internal server error")
	}

	// Parse domain.PresignedURL to model.PresignedURL
	urls := make([]*model.PresignedURL, len(res.Urls))
	for i, url := range res.Urls {
		urls[i] = &model.PresignedURL{
			ObjectID: url.ObjectId,
			URL:      url.Url,
			File:     url.File,
		}
	}

	return &model.CreatePresignedUrlsResponse{
		Urls: urls,
	}, nil
}

// SoftDeleteNote is the resolver for the softDeleteNotes field.
func SoftDeleteNote(ctx context.Context, id string, srv service.NoteService) (bool, error) {
	// Check if the user is authenticated
	userId := domain.GetUserIdFromContext(ctx)
	if userId == uuid.Nil {
		return false, errors.New("unauthenticated")
	}

	noteId, err := uuid.Parse(id)
	if err != nil {
		return false, errors.New("invalid note id")
	}

	err = srv.DeleteNote(ctx, noteId, false)
	if err != nil {
		switch err.Error() {
		case "note not found":
			return false, errors.New("note not found")
		default:
			return false, errors.New("internal server error")
		}
	}

	return true, nil
}

// DeleteNote is the resolver for the deleteNote field.
func DeleteNote(ctx context.Context, id string, srv service.NoteService) (bool, error) {
	// Check if the user is authenticated
	userId := domain.GetUserIdFromContext(ctx)
	if userId == uuid.Nil {
		return false, errors.New("unauthenticated")
	}

	noteId, err := uuid.Parse(id)
	if err != nil {
		return false, errors.New("invalid note id")
	}

	err = srv.DeleteNote(ctx, noteId, true)
	if err != nil {
		switch err.Error() {
		case "note not found":
			return false, errors.New("note not found")
		default:
			return false, errors.New("internal server error")
		}
	}

	return true, nil
}
