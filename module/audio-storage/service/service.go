package service

import (
	"context"

	"github.com/dibikhairurrazi/audio-storage/module/audio-storage/model"
)

type UserService interface {
	FindUser(context.Context, int) (model.User, error)
}

type PhraseService interface {
	Store(context.Context, model.Phrase) error
	Retrieve(context.Context, int, int, string) (model.Phrase, error)
}
