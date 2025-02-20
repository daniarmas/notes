// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/daniarmas/notes/internal/domain"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// AccessTokenRepository is an autogenerated mock type for the AccessTokenRepository type
type AccessTokenRepository struct {
	mock.Mock
}

// CreateAccessToken provides a mock function with given fields: ctx, userId, refreshTokenId
func (_m *AccessTokenRepository) CreateAccessToken(ctx context.Context, userId uuid.UUID, refreshTokenId uuid.UUID) (*domain.AccessToken, error) {
	ret := _m.Called(ctx, userId, refreshTokenId)

	if len(ret) == 0 {
		panic("no return value specified for CreateAccessToken")
	}

	var r0 *domain.AccessToken
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) (*domain.AccessToken, error)); ok {
		return rf(ctx, userId, refreshTokenId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) *domain.AccessToken); ok {
		r0 = rf(ctx, userId, refreshTokenId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.AccessToken)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(ctx, userId, refreshTokenId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAccessTokenByUserId provides a mock function with given fields: ctx, userId
func (_m *AccessTokenRepository) DeleteAccessTokenByUserId(ctx context.Context, userId uuid.UUID) error {
	ret := _m.Called(ctx, userId)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAccessTokenByUserId")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, userId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAccessToken provides a mock function with given fields: ctx, id
func (_m *AccessTokenRepository) GetAccessToken(ctx context.Context, id uuid.UUID) (*domain.AccessToken, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetAccessToken")
	}

	var r0 *domain.AccessToken
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*domain.AccessToken, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *domain.AccessToken); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.AccessToken)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAccessTokenRepository creates a new instance of AccessTokenRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAccessTokenRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *AccessTokenRepository {
	mock := &AccessTokenRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
