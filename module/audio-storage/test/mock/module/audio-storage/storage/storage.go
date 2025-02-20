// Code generated by MockGen. DO NOT EDIT.
// Source: ./module/audio-storage/storage/storage.go
//
// Generated by this command:
//
//	mockgen -source ./module/audio-storage/storage/storage.go -destination module/audio-storage/test/mock/./module/audio-storage/storage/storage.go
//

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	context "context"
	reflect "reflect"

	model "github.com/dibikhairurrazi/audio-storage/module/audio-storage/model"
	gomock "go.uber.org/mock/gomock"
)

// MockAudioStorage is a mock of AudioStorage interface.
type MockAudioStorage struct {
	ctrl     *gomock.Controller
	recorder *MockAudioStorageMockRecorder
}

// MockAudioStorageMockRecorder is the mock recorder for MockAudioStorage.
type MockAudioStorageMockRecorder struct {
	mock *MockAudioStorage
}

// NewMockAudioStorage creates a new mock instance.
func NewMockAudioStorage(ctrl *gomock.Controller) *MockAudioStorage {
	mock := &MockAudioStorage{ctrl: ctrl}
	mock.recorder = &MockAudioStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAudioStorage) EXPECT() *MockAudioStorageMockRecorder {
	return m.recorder
}

// DeleteFile mocks base method.
func (m *MockAudioStorage) DeleteFile(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFile", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFile indicates an expected call of DeleteFile.
func (mr *MockAudioStorageMockRecorder) DeleteFile(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFile", reflect.TypeOf((*MockAudioStorage)(nil).DeleteFile), arg0, arg1)
}

// LoadFile mocks base method.
func (m *MockAudioStorage) LoadFile(arg0 context.Context, arg1 string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadFile", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadFile indicates an expected call of LoadFile.
func (mr *MockAudioStorageMockRecorder) LoadFile(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadFile", reflect.TypeOf((*MockAudioStorage)(nil).LoadFile), arg0, arg1)
}

// SaveFile mocks base method.
func (m *MockAudioStorage) SaveFile(arg0 context.Context, arg1, arg2 string, arg3 model.Phrase) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveFile", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveFile indicates an expected call of SaveFile.
func (mr *MockAudioStorageMockRecorder) SaveFile(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveFile", reflect.TypeOf((*MockAudioStorage)(nil).SaveFile), arg0, arg1, arg2, arg3)
}
