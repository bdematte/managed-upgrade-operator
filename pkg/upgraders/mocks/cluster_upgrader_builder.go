// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/openshift/managed-upgrade-operator/pkg/upgraders (interfaces: ClusterUpgraderBuilder)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	v1alpha1 "github.com/openshift/managed-upgrade-operator/api/v1alpha1"
	configmanager "github.com/openshift/managed-upgrade-operator/pkg/configmanager"
	eventmanager "github.com/openshift/managed-upgrade-operator/pkg/eventmanager"
	metrics "github.com/openshift/managed-upgrade-operator/pkg/metrics"
	upgraders "github.com/openshift/managed-upgrade-operator/pkg/upgraders"
	reflect "reflect"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

// MockClusterUpgraderBuilder is a mock of ClusterUpgraderBuilder interface
type MockClusterUpgraderBuilder struct {
	ctrl     *gomock.Controller
	recorder *MockClusterUpgraderBuilderMockRecorder
}

// MockClusterUpgraderBuilderMockRecorder is the mock recorder for MockClusterUpgraderBuilder
type MockClusterUpgraderBuilderMockRecorder struct {
	mock *MockClusterUpgraderBuilder
}

// NewMockClusterUpgraderBuilder creates a new mock instance
func NewMockClusterUpgraderBuilder(ctrl *gomock.Controller) *MockClusterUpgraderBuilder {
	mock := &MockClusterUpgraderBuilder{ctrl: ctrl}
	mock.recorder = &MockClusterUpgraderBuilderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClusterUpgraderBuilder) EXPECT() *MockClusterUpgraderBuilderMockRecorder {
	return m.recorder
}

// NewClient mocks base method
func (m *MockClusterUpgraderBuilder) NewClient(arg0 client.Client, arg1 configmanager.ConfigManager, arg2 metrics.Metrics, arg3 eventmanager.EventManager, arg4 v1alpha1.UpgradeType) (upgraders.ClusterUpgrader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewClient", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(upgraders.ClusterUpgrader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewClient indicates an expected call of NewClient
func (mr *MockClusterUpgraderBuilderMockRecorder) NewClient(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewClient", reflect.TypeOf((*MockClusterUpgraderBuilder)(nil).NewClient), arg0, arg1, arg2, arg3, arg4)
}
