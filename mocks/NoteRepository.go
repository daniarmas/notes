// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/daniarmas/notes/internal/domain"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// NoteRepository is an autogenerated mock type for the NoteRepository type
type NoteRepository struct {
	mock.Mock
}

// CreateNote provides a mock function with given fields: ctx, note
func (_m *NoteRepository) CreateNote(ctx context.Context, note *domain.Note) (*domain.Note, error) {
	ret := _m.Called(ctx, note)

	if len(ret) == 0 {
		panic("no return value specified for CreateNote")
	}

	var r0 *domain.Note
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Note) (*domain.Note, error)); ok {
		return rf(ctx, note)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Note) *domain.Note); ok {
		r0 = rf(ctx, note)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Note)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *domain.Note) error); ok {
		r1 = rf(ctx, note)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteNote provides a mock function with given fields: ctx, id
func (_m *NoteRepository) DeleteNote(ctx context.Context, id uuid.UUID) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteNote")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetNote provides a mock function with given fields: ctx, id
func (_m *NoteRepository) GetNote(ctx context.Context, id uuid.UUID) (*domain.Note, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetNote")
	}

	var r0 *domain.Note
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*domain.Note, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *domain.Note); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Note)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListNote provides a mock function with given fields: ctx
func (_m *NoteRepository) ListNote(ctx context.Context) (*[]domain.Note, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for ListNote")
	}

	var r0 *[]domain.Note
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*[]domain.Note, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *[]domain.Note); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]domain.Note)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateNote provides a mock function with given fields: ctx, note
func (_m *NoteRepository) UpdateNote(ctx context.Context, note *domain.Note) (*domain.Note, error) {
	ret := _m.Called(ctx, note)

	if len(ret) == 0 {
		panic("no return value specified for UpdateNote")
	}

	var r0 *domain.Note
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Note) (*domain.Note, error)); ok {
		return rf(ctx, note)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Note) *domain.Note); ok {
		r0 = rf(ctx, note)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Note)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *domain.Note) error); ok {
		r1 = rf(ctx, note)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewNoteRepository creates a new instance of NoteRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewNoteRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *NoteRepository {
	mock := &NoteRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}