// Code generated by MockGen. DO NOT EDIT.
// Source: ./module/audio-storage/service/service.go
//
// Generated by this command:
//
//	mockgen -source ./module/audio-storage/service/service.go -destination module/audio-storage/test/mock/./module/audio-storage/service/service.go
//

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	model "github.com/dibikhairurrazi/audio-storage/module/audio-storage/model"
	gomock "go.uber.org/mock/gomock"
)

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// FindUser mocks base method.
func (m *MockUserService) FindUser(arg0 context.Context, arg1 int) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUser", arg0, arg1)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUser indicates an expected call of FindUser.
func (mr *MockUserServiceMockRecorder) FindUser(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUser", reflect.TypeOf((*MockUserService)(nil).FindUser), arg0, arg1)
}

// MockPhraseService is a mock of PhraseService interface.
type MockPhraseService struct {
	ctrl     *gomock.Controller
	recorder *MockPhraseServiceMockRecorder
}

// MockPhraseServiceMockRecorder is the mock recorder for MockPhraseService.
type MockPhraseServiceMockRecorder struct {
	mock *MockPhraseService
}

// NewMockPhraseService creates a new mock instance.
func NewMockPhraseService(ctrl *gomock.Controller) *MockPhraseService {
	mock := &MockPhraseService{ctrl: ctrl}
	mock.recorder = &MockPhraseServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPhraseService) EXPECT() *MockPhraseServiceMockRecorder {
	return m.recorder
}

// Retrieve mocks base method.
func (m *MockPhraseService) Retrieve(arg0 context.Context, arg1, arg2 int, arg3 string) (model.Phrase, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Retrieve", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(model.Phrase)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Retrieve indicates an expected call of Retrieve.
func (mr *MockPhraseServiceMockRecorder) Retrieve(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Retrieve", reflect.TypeOf((*MockPhraseService)(nil).Retrieve), arg0, arg1, arg2, arg3)
}

// Store mocks base method.
func (m *MockPhraseService) Store(arg0 context.Context, arg1 model.Phrase) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Store indicates an expected call of Store.
func (mr *MockPhraseServiceMockRecorder) Store(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockPhraseService)(nil).Store), arg0, arg1)
}
