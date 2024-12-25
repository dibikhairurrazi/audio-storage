package service

import (
	"context"
	"log/slog"

	"github.com/dibikhairurrazi/audio-storage/module/audio-storage/model"
	"github.com/dibikhairurrazi/audio-storage/module/audio-storage/repository"
)

type UserServiceProvider struct {
	UserRepo repository.UserRepository
}

func NewUserServiceProvider(userRepo repository.UserRepository) *UserServiceProvider {
	return &UserServiceProvider{
		UserRepo: userRepo,
	}
}

func (us *UserServiceProvider) FindUser(ctx context.Context, userID int) (model.User, error) {
	user, err := us.UserRepo.FindUser(ctx, userID)
	if err != nil {
		slog.Error("error finding user", "user ID", userID)
	}

	return user, nil
}
