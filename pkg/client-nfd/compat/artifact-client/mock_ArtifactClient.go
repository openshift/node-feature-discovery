// Code generated by mockery v2.42.0. DO NOT EDIT.

package compat

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	v1alpha1 "github.com/openshift/node-feature-discovery/api/image-compatibility/v1alpha1"
)

// MockArtifactClient is an autogenerated mock type for the ArtifactClient type
type MockArtifactClient struct {
	mock.Mock
}

// FetchCompatibilitySpec provides a mock function with given fields: ctx
func (_m *MockArtifactClient) FetchCompatibilitySpec(ctx context.Context) (*v1alpha1.Spec, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for FetchCompatibilitySpec")
	}

	var r0 *v1alpha1.Spec
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*v1alpha1.Spec, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *v1alpha1.Spec); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1alpha1.Spec)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockArtifactClient creates a new instance of MockArtifactClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockArtifactClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockArtifactClient {
	mock := &MockArtifactClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
