package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.63

import (
	"context"

	"github.com/daniarmas/notes/internal/graph/model"
	"github.com/daniarmas/notes/internal/graph/resolver"
)

// SignIn is the resolver for the signIn field.
func (r *mutationResolver) SignIn(ctx context.Context, input model.SignInInput) (*model.SignInResponse, error) {
	return resolver.SignIn(ctx, input, r.AuthSrv)
}

// SignOut is the resolver for the signOut field.
func (r *mutationResolver) SignOut(ctx context.Context) (bool, error) {
	return resolver.SignOut(ctx, r.AuthSrv)
}

// CreateNote is the resolver for the createNote field.
func (r *mutationResolver) CreateNote(ctx context.Context, input model.CreateNoteInput) (*model.Note, error) {
	return resolver.CreateNote(ctx, input, r.NoteSrv)
}

// CreatePresignedURL is the resolver for the createPresignedUrl field.
func (r *mutationResolver) CreatePresignedURL(ctx context.Context, objectName []string) (*model.CreatePresignedUrlsResponse, error) {
	return resolver.CreatePresignedURL(ctx, objectName, r.NoteSrv)
}

// SoftDeleteNote is the resolver for the softDeleteNote field.
func (r *mutationResolver) SoftDeleteNote(ctx context.Context, id string) (bool, error) {
	return resolver.SoftDeleteNote(ctx, id, r.NoteSrv)
}

// DeleteNote is the resolver for the deleteNote field.
func (r *mutationResolver) DeleteNote(ctx context.Context, id string) (bool, error) {
	return resolver.DeleteNote(ctx, id, r.NoteSrv)
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	return resolver.Me(ctx, r.AuthSrv)
}

// ListNotes is the resolver for the listNotes field.
func (r *queryResolver) ListNotes(ctx context.Context, input *model.NotesInput) (*model.NotesResponse, error) {
	return resolver.ListNotes(ctx, input, r.NoteSrv)
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
/*
	func (r *mutationResolver) SoftDeleteNotes(ctx context.Context, id string) (bool, error) {
	return resolver.SoftDeleteNotes(ctx, id, r.NoteSrv)
}
func (r *mutationResolver) DeleteNotes(ctx context.Context, id string) (bool, error) {
	panic(fmt.Errorf("not implemented: DeleteNotes - deleteNotes"))
}
*/
