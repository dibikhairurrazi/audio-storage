package repository

import (
	"context"
	"database/sql"

	"github.com/dibikhairurrazi/audio-storage/module/audio-storage/model"
)

type UserRepository interface {
	FindUser(context.Context, int) (model.User, error)
}

type PhraseRepository interface {
	CreateTx(context.Context) (*sql.Tx, error)
	SavePhraseMetadata(context.Context, model.Phrase) error
	FindPhraseMetadata(context.Context, int, int) (model.Phrase, error)
}
