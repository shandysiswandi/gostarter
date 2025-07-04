// Code generated by mockery. DO NOT EDIT.

package mockz

import (
	context "context"

	enum "github.com/shandysiswandi/goreng/enum"
	domain "github.com/shandysiswandi/gostarter/internal/todo/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// MockUpdateStatusStore is an autogenerated mock type for the UpdateStatusStore type
type MockUpdateStatusStore struct {
	mock.Mock
}

type MockUpdateStatusStore_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUpdateStatusStore) EXPECT() *MockUpdateStatusStore_Expecter {
	return &MockUpdateStatusStore_Expecter{mock: &_m.Mock}
}

// UpdateStatus provides a mock function with given fields: ctx, in, sts
func (_m *MockUpdateStatusStore) UpdateStatus(ctx context.Context, in uint64, sts enum.Enum[domain.TodoStatus]) error {
	ret := _m.Called(ctx, in, sts)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStatus")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, enum.Enum[domain.TodoStatus]) error); ok {
		r0 = rf(ctx, in, sts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockUpdateStatusStore_UpdateStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateStatus'
type MockUpdateStatusStore_UpdateStatus_Call struct {
	*mock.Call
}

// UpdateStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - in uint64
//   - sts enum.Enum[domain.TodoStatus]
func (_e *MockUpdateStatusStore_Expecter) UpdateStatus(ctx interface{}, in interface{}, sts interface{}) *MockUpdateStatusStore_UpdateStatus_Call {
	return &MockUpdateStatusStore_UpdateStatus_Call{Call: _e.mock.On("UpdateStatus", ctx, in, sts)}
}

func (_c *MockUpdateStatusStore_UpdateStatus_Call) Run(run func(ctx context.Context, in uint64, sts enum.Enum[domain.TodoStatus])) *MockUpdateStatusStore_UpdateStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(enum.Enum[domain.TodoStatus]))
	})
	return _c
}

func (_c *MockUpdateStatusStore_UpdateStatus_Call) Return(_a0 error) *MockUpdateStatusStore_UpdateStatus_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockUpdateStatusStore_UpdateStatus_Call) RunAndReturn(run func(context.Context, uint64, enum.Enum[domain.TodoStatus]) error) *MockUpdateStatusStore_UpdateStatus_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockUpdateStatusStore creates a new instance of MockUpdateStatusStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUpdateStatusStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUpdateStatusStore {
	mock := &MockUpdateStatusStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
