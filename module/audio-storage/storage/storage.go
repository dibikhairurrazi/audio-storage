package storage

import (
	"context"

	"github.com/dibikhairurrazi/audio-storage/module/audio-storage/model"
)

type AudioStorage interface {
	SaveFile(context.Context, string, string, model.Phrase) (string, error)
	LoadFile(context.Context, string) ([]byte, error)
	DeleteFile(context.Context, string) error
}
