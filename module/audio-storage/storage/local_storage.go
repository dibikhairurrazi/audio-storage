package storage

import (
	"context"
	"os"
	"path"

	"github.com/dibikhairurrazi/audio-storage/module/audio-storage/model"
)

type LocalStorage struct {
	RootFolder string
}

func NewLocalStorage(rootFolder string) *LocalStorage {
	return &LocalStorage{RootFolder: rootFolder}
}

func (ls *LocalStorage) SaveFile(ctx context.Context, folder string, filename string, phrase model.Phrase) (string, error) {
	// create upload folder is not exists
	if _, err := os.Stat(ls.RootFolder); os.IsNotExist(err) {
		os.Mkdir(ls.RootFolder, 0755)
	}

	// create user folder if not exists
	if _, err := os.Stat(path.Join(ls.RootFolder, folder)); os.IsNotExist(err) {
		os.Mkdir(path.Join(ls.RootFolder, folder), 0755)
	}

	filePath := path.Join(ls.RootFolder, folder, filename)
	newFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer newFile.Close()
	if _, err := newFile.Write(phrase.Content); err != nil {
		return "", err
	}

	return filePath, nil
}

func (ls *LocalStorage) LoadFile(ctx context.Context, filepath string) ([]byte, error) {
	return os.ReadFile(filepath)
}

func (ls *LocalStorage) DeleteFile(ctx context.Context, filepath string) error {
	return os.Remove(filepath)
}
