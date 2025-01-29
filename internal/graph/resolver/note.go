package resolver

import (
	"context"
	"errors"
	"time"

	"github.com/daniarmas/notes/internal/domain"
	"github.com/daniarmas/notes/internal/graph/model"
	"github.com/daniarmas/notes/internal/service"
	"github.com/daniarmas/notes/internal/utils"
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

// Notes is the resolver for the notes field.
func Notes(ctx context.Context, input *model.NotesInput, srv service.NoteService) (*model.NotesResponse, error) {
	// Get the cursor from the query parameters
	var cursorQueryParam string
	if input != nil {
		cursorQueryParam = input.Cursor
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

	notes, err := srv.ListNotesByUser(ctx, cursor)
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
