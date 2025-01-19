// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	sql "database/sql"
)

// DBTX is an autogenerated mock type for the DBTX type
type DBTX struct {
	mock.Mock
}

// ExecContext provides a mock function with given fields: _a0, _a1, _a2
func (_m *DBTX) ExecContext(_a0 context.Context, _a1 string, _a2 ...interface{}) (sql.Result, error) {
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _a2...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ExecContext")
	}

	var r0 sql.Result
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) (sql.Result, error)); ok {
		return rf(_a0, _a1, _a2...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) sql.Result); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sql.Result)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, ...interface{}) error); ok {
		r1 = rf(_a0, _a1, _a2...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PrepareContext provides a mock function with given fields: _a0, _a1
func (_m *DBTX) PrepareContext(_a0 context.Context, _a1 string) (*sql.Stmt, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for PrepareContext")
	}

	var r0 *sql.Stmt
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*sql.Stmt, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *sql.Stmt); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Stmt)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryContext provides a mock function with given fields: _a0, _a1, _a2
func (_m *DBTX) QueryContext(_a0 context.Context, _a1 string, _a2 ...interface{}) (*sql.Rows, error) {
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _a2...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for QueryContext")
	}

	var r0 *sql.Rows
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) (*sql.Rows, error)); ok {
		return rf(_a0, _a1, _a2...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) *sql.Rows); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Rows)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, ...interface{}) error); ok {
		r1 = rf(_a0, _a1, _a2...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryRowContext provides a mock function with given fields: _a0, _a1, _a2
func (_m *DBTX) QueryRowContext(_a0 context.Context, _a1 string, _a2 ...interface{}) *sql.Row {
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _a2...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for QueryRowContext")
	}

	var r0 *sql.Row
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) *sql.Row); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Row)
		}
	}

	return r0
}

// NewDBTX creates a new instance of DBTX. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDBTX(t interface {
	mock.TestingT
	Cleanup(func())
}) *DBTX {
	mock := &DBTX{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
