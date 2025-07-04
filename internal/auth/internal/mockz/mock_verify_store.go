// Code generated by mockery. DO NOT EDIT.

package mockz

import (
	context "context"

	domain "github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// MockVerifyStore is an autogenerated mock type for the VerifyStore type
type MockVerifyStore struct {
	mock.Mock
}

type MockVerifyStore_Expecter struct {
	mock *mock.Mock
}

func (_m *MockVerifyStore) EXPECT() *MockVerifyStore_Expecter {
	return &MockVerifyStore_Expecter{mock: &_m.Mock}
}

// UserByEmail provides a mock function with given fields: ctx, email
func (_m *MockVerifyStore) UserByEmail(ctx context.Context, email string) (*domain.User, error) {
	ret := _m.Called(ctx, email)

	if len(ret) == 0 {
		panic("no return value specified for UserByEmail")
	}

	var r0 *domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*domain.User, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.User); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockVerifyStore_UserByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UserByEmail'
type MockVerifyStore_UserByEmail_Call struct {
	*mock.Call
}

// UserByEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - email string
func (_e *MockVerifyStore_Expecter) UserByEmail(ctx interface{}, email interface{}) *MockVerifyStore_UserByEmail_Call {
	return &MockVerifyStore_UserByEmail_Call{Call: _e.mock.On("UserByEmail", ctx, email)}
}

func (_c *MockVerifyStore_UserByEmail_Call) Run(run func(ctx context.Context, email string)) *MockVerifyStore_UserByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockVerifyStore_UserByEmail_Call) Return(_a0 *domain.User, _a1 error) *MockVerifyStore_UserByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockVerifyStore_UserByEmail_Call) RunAndReturn(run func(context.Context, string) (*domain.User, error)) *MockVerifyStore_UserByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// UserVerificationByUserID provides a mock function with given fields: ctx, uid
func (_m *MockVerifyStore) UserVerificationByUserID(ctx context.Context, uid uint64) (*domain.UserVerification, error) {
	ret := _m.Called(ctx, uid)

	if len(ret) == 0 {
		panic("no return value specified for UserVerificationByUserID")
	}

	var r0 *domain.UserVerification
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) (*domain.UserVerification, error)); ok {
		return rf(ctx, uid)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64) *domain.UserVerification); ok {
		r0 = rf(ctx, uid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.UserVerification)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, uid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockVerifyStore_UserVerificationByUserID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UserVerificationByUserID'
type MockVerifyStore_UserVerificationByUserID_Call struct {
	*mock.Call
}

// UserVerificationByUserID is a helper method to define mock.On call
//   - ctx context.Context
//   - uid uint64
func (_e *MockVerifyStore_Expecter) UserVerificationByUserID(ctx interface{}, uid interface{}) *MockVerifyStore_UserVerificationByUserID_Call {
	return &MockVerifyStore_UserVerificationByUserID_Call{Call: _e.mock.On("UserVerificationByUserID", ctx, uid)}
}

func (_c *MockVerifyStore_UserVerificationByUserID_Call) Run(run func(ctx context.Context, uid uint64)) *MockVerifyStore_UserVerificationByUserID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64))
	})
	return _c
}

func (_c *MockVerifyStore_UserVerificationByUserID_Call) Return(_a0 *domain.UserVerification, _a1 error) *MockVerifyStore_UserVerificationByUserID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockVerifyStore_UserVerificationByUserID_Call) RunAndReturn(run func(context.Context, uint64) (*domain.UserVerification, error)) *MockVerifyStore_UserVerificationByUserID_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockVerifyStore creates a new instance of MockVerifyStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockVerifyStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockVerifyStore {
	mock := &MockVerifyStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
