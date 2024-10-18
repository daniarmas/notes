// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// HashDatasource is an autogenerated mock type for the HashDatasource type
type HashDatasource struct {
	mock.Mock
}

// CheckHash provides a mock function with given fields: value, hash
func (_m *HashDatasource) CheckHash(value string, hash string) (bool, error) {
	ret := _m.Called(value, hash)

	if len(ret) == 0 {
		panic("no return value specified for CheckHash")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (bool, error)); ok {
		return rf(value, hash)
	}
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(value, hash)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(value, hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Hash provides a mock function with given fields: value
func (_m *HashDatasource) Hash(value string) (string, error) {
	ret := _m.Called(value)

	if len(ret) == 0 {
		panic("no return value specified for Hash")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(value)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(value)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewHashDatasource creates a new instance of HashDatasource. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHashDatasource(t interface {
	mock.TestingT
	Cleanup(func())
}) *HashDatasource {
	mock := &HashDatasource{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}