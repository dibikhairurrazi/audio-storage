package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/dibikhairurrazi/audio-storage/module/audio-storage/model"
	mock_repository "github.com/dibikhairurrazi/audio-storage/module/audio-storage/test/mock/module/audio-storage/repository"
	"go.uber.org/mock/gomock"
)

var (
	testUser = model.User{
		ID:    defaultUserID,
		Name:  "John Doe",
		Email: "john.doe@something.com",
	}
)

func TestNewUserServiceProvider(t *testing.T) {
	type args struct {
		userRepo *mock_repository.MockUserRepository
	}
	tests := []struct {
		name string
		args args
		want *UserServiceProvider
	}{
		{
			name: "constructor",
			args: args{
				userRepo: mock_repository.NewMockUserRepository(gomock.NewController(t)),
			},
			want: &UserServiceProvider{
				UserRepo: mock_repository.NewMockUserRepository(gomock.NewController(t)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserServiceProvider(tt.args.userRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserServiceProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserServiceProvider_FindUser(t *testing.T) {
	type fields struct {
		UserRepo *mock_repository.MockUserRepository
	}
	type args struct {
		ctx    context.Context
		userID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.User
		wantErr bool
		prepare func(f *fields)
	}{
		{
			name: "happy path",
			fields: fields{
				UserRepo: mock_repository.NewMockUserRepository(gomock.NewController(t)),
			},
			args: args{
				ctx:    context.Background(),
				userID: defaultUserID,
			},
			want:    testUser,
			wantErr: false,
			prepare: func(f *fields) {
				f.UserRepo.EXPECT().FindUser(context.Background(), defaultUserID).Return(testUser, nil)
			},
		},
		{
			name: "error finding user",
			fields: fields{
				UserRepo: mock_repository.NewMockUserRepository(gomock.NewController(t)),
			},
			args: args{
				ctx:    context.Background(),
				userID: defaultUserID,
			},
			want:    model.User{},
			wantErr: true,
			prepare: func(f *fields) {
				f.UserRepo.EXPECT().FindUser(context.Background(), defaultUserID).Return(model.User{}, errors.New("db error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserServiceProvider{
				UserRepo: tt.fields.UserRepo,
			}
			tt.prepare(&tt.fields)
			got, err := us.FindUser(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserServiceProvider.FindUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserServiceProvider.FindUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
