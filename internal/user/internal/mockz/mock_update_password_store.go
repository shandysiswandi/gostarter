// Code generated by mockery. DO NOT EDIT.

package mockz

import (
	context "context"

	domain "github.com/shandysiswandi/gostarter/internal/user/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// MockUpdatePasswordStore is an autogenerated mock type for the UpdatePasswordStore type
type MockUpdatePasswordStore struct {
	mock.Mock
}

type MockUpdatePasswordStore_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUpdatePasswordStore) EXPECT() *MockUpdatePasswordStore_Expecter {
	return &MockUpdatePasswordStore_Expecter{mock: &_m.Mock}
}

// FindUser provides a mock function with given fields: ctx, id
func (_m *MockUpdatePasswordStore) FindUser(ctx context.Context, id uint64) (*domain.User, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for FindUser")
	}

	var r0 *domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) (*domain.User, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64) *domain.User); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUpdatePasswordStore_FindUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindUser'
type MockUpdatePasswordStore_FindUser_Call struct {
	*mock.Call
}

// FindUser is a helper method to define mock.On call
//   - ctx context.Context
//   - id uint64
func (_e *MockUpdatePasswordStore_Expecter) FindUser(ctx interface{}, id interface{}) *MockUpdatePasswordStore_FindUser_Call {
	return &MockUpdatePasswordStore_FindUser_Call{Call: _e.mock.On("FindUser", ctx, id)}
}

func (_c *MockUpdatePasswordStore_FindUser_Call) Run(run func(ctx context.Context, id uint64)) *MockUpdatePasswordStore_FindUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64))
	})
	return _c
}

func (_c *MockUpdatePasswordStore_FindUser_Call) Return(_a0 *domain.User, _a1 error) *MockUpdatePasswordStore_FindUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUpdatePasswordStore_FindUser_Call) RunAndReturn(run func(context.Context, uint64) (*domain.User, error)) *MockUpdatePasswordStore_FindUser_Call {
	_c.Call.Return(run)
	return _c
}

// UpdatePassword provides a mock function with given fields: ctx, user
func (_m *MockUpdatePasswordStore) UpdatePassword(ctx context.Context, user domain.User) error {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for UpdatePassword")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockUpdatePasswordStore_UpdatePassword_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdatePassword'
type MockUpdatePasswordStore_UpdatePassword_Call struct {
	*mock.Call
}

// UpdatePassword is a helper method to define mock.On call
//   - ctx context.Context
//   - user domain.User
func (_e *MockUpdatePasswordStore_Expecter) UpdatePassword(ctx interface{}, user interface{}) *MockUpdatePasswordStore_UpdatePassword_Call {
	return &MockUpdatePasswordStore_UpdatePassword_Call{Call: _e.mock.On("UpdatePassword", ctx, user)}
}

func (_c *MockUpdatePasswordStore_UpdatePassword_Call) Run(run func(ctx context.Context, user domain.User)) *MockUpdatePasswordStore_UpdatePassword_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.User))
	})
	return _c
}

func (_c *MockUpdatePasswordStore_UpdatePassword_Call) Return(_a0 error) *MockUpdatePasswordStore_UpdatePassword_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockUpdatePasswordStore_UpdatePassword_Call) RunAndReturn(run func(context.Context, domain.User) error) *MockUpdatePasswordStore_UpdatePassword_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockUpdatePasswordStore creates a new instance of MockUpdatePasswordStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUpdatePasswordStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUpdatePasswordStore {
	mock := &MockUpdatePasswordStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
