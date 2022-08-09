// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	health "github.com/Bhinneka/user-service/protogo/health"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	mock "github.com/stretchr/testify/mock"
)

// PingPongServiceServer is an autogenerated mock type for the PingPongServiceServer type
type PingPongServiceServer struct {
	mock.Mock
}

// PingPong provides a mock function with given fields: _a0, _a1
func (_m *PingPongServiceServer) PingPong(_a0 context.Context, _a1 *emptypb.Empty) (*health.PongResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *health.PongResponse
	if rf, ok := ret.Get(0).(func(context.Context, *emptypb.Empty) *health.PongResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*health.PongResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *emptypb.Empty) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
