package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path"
	"strconv"
	"time"

	"github.com/dibikhairurrazi/audio-storage/module/audio-storage/converter"
	"github.com/dibikhairurrazi/audio-storage/module/audio-storage/model"
	"github.com/dibikhairurrazi/audio-storage/module/audio-storage/repository"
	"github.com/dibikhairurrazi/audio-storage/module/audio-storage/storage"
	"golang.org/x/exp/rand"
)

const (
	defaultFilenameLength = 10
	defaultSavedFormat    = "wav"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type PhraseServiceProvider struct {
	Repository repository.PhraseRepository
	Storage    storage.AudioStorage
	Converter  converter.Converter
}

func NewPhraseServiceProvider(phraseRepo repository.PhraseRepository, audioConverter converter.Converter, storage storage.AudioStorage) *PhraseServiceProvider {
	return &PhraseServiceProvider{
		Repository: phraseRepo,
		Converter:  audioConverter,
		Storage:    storage,
	}
}

func (ps *PhraseServiceProvider) Store(ctx context.Context, phrase model.Phrase) error {
	now := time.Now()
	defer func() {
		slog.Info(fmt.Sprintf("success storing file, the process took %v ms", time.Since(now).Milliseconds()))
	}()

	converted, err := ps.Converter.Convert(phrase.Content, defaultSavedFormat)
	if err != nil {
		slog.Error("error converting file", "userID", phrase.UserID, "err", err.Error())
		return err
	}

	phrase.Content = converted
	savedPath, err := ps.Storage.SaveFile(ctx, path.Join(strconv.Itoa(phrase.UserID), strconv.Itoa(phrase.ID)), fmt.Sprintf("%v.%v", generateRandomStringwithNLength(defaultFilenameLength), defaultSavedFormat), phrase)
	if err != nil {
		slog.Error("error saving file to storage", "userID", phrase.UserID, "err", err.Error())
		return err
	}

	phrase.FilePath = savedPath
	err = ps.Repository.SavePhraseMetadata(ctx, phrase)
	if err != nil {
		slog.Error("failed to save phrase metadata to DB, deleting saved file", "userID", phrase.UserID, "err", err.Error())

		err2 := ps.Storage.DeleteFile(ctx, savedPath)
		if err2 != nil {
			slog.Error("failed to delete file", "path", savedPath, "err", err.Error())
		}

		return err
	}

	return nil
}

func (ps *PhraseServiceProvider) Retrieve(ctx context.Context, userID, phraseID int, extension string) (model.Phrase, error) {
	now := time.Now()
	defer func() {
		slog.Info(fmt.Sprintf("success retrieving file, the process took %v ms", time.Since(now).Milliseconds()))
	}()

	phrase, err := ps.Repository.FindPhraseMetadata(ctx, userID, phraseID)
	if err != nil {
		slog.Error("failed to find phrase meta data", "userID", userID, "phraseID", phraseID, "err", err.Error())
		return model.Phrase{}, err
	}

	if phrase.FilePath == "" {
		slog.Error("phrase contain empty filepath", "phraseID", phraseID)
		return model.Phrase{}, errors.New("invalid filepath")
	}

	content, err := ps.Storage.LoadFile(ctx, phrase.FilePath)
	if err != nil {
		slog.Error("failed to find open file content", "userID", userID, "phraseID", phraseID, "err", err.Error())
		return model.Phrase{}, err
	}

	convertedContent, err := ps.Converter.Convert(content, extension)
	if err != nil {
		slog.Error("failed to convert audio file", "userID", userID, "phraseID", phraseID, "err", err.Error())
		return model.Phrase{}, err
	}

	phrase.Content = convertedContent

	return phrase, nil
}

func generateRandomStringwithNLength(n int) string {
	b := make([]rune, n)
	rand.Seed(uint64(time.Now().UnixNano()))
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
