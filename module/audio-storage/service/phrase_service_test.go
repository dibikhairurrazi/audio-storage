package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/dibikhairurrazi/audio-storage/module/audio-storage/model"
	mock_converter "github.com/dibikhairurrazi/audio-storage/module/audio-storage/test/mock/module/audio-storage/converter"
	mock_repository "github.com/dibikhairurrazi/audio-storage/module/audio-storage/test/mock/module/audio-storage/repository"
	mock_storage "github.com/dibikhairurrazi/audio-storage/module/audio-storage/test/mock/module/audio-storage/storage"
	"go.uber.org/mock/gomock"
)

var (
	defaultUserID      = 1
	defaultPhraseID    = 1
	defaultUserIDinStr = "1"

	testContent = prepareTestAudioContent()
	testPhrase  = model.Phrase{
		ID:      defaultUserID,
		UserID:  defaultPhraseID,
		Content: testContent,
	}

	savedPath          = "upload/1/s012oi1dmn.wav"
	testPhraseWithPath = model.Phrase{
		ID:       defaultUserID,
		UserID:   defaultPhraseID,
		Content:  testContent,
		FilePath: savedPath,
	}

	ext = "mp3"
)

func TestPhraseServiceProvider_Store(t *testing.T) {
	type fields struct {
		Repository *mock_repository.MockPhraseRepository
		Storage    *mock_storage.MockAudioStorage
		Converter  *mock_converter.MockConverter
	}
	type args struct {
		ctx    context.Context
		phrase model.Phrase
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		prepare func(f *fields)
	}{
		{
			name: "happy path",
			fields: fields{
				Repository: mock_repository.NewMockPhraseRepository(gomock.NewController(t)),
				Storage:    mock_storage.NewMockAudioStorage(gomock.NewController(t)),
				Converter:  mock_converter.NewMockConverter(gomock.NewController(t)),
			},
			args: args{
				ctx:    context.Background(),
				phrase: testPhrase,
			},
			wantErr: false,
			prepare: func(f *fields) {
				f.Converter.EXPECT().Convert(testContent, "wav").Return(testContent, nil)
				f.Storage.EXPECT().SaveFile(context.Background(), defaultUserIDinStr, gomock.Any(), testPhrase).Return(savedPath, nil)
				f.Repository.EXPECT().SavePhraseMetadata(context.Background(), testPhraseWithPath).Return(nil)
			},
		},
		{
			name: "failed to convert audio files",
			fields: fields{
				Repository: mock_repository.NewMockPhraseRepository(gomock.NewController(t)),
				Storage:    mock_storage.NewMockAudioStorage(gomock.NewController(t)),
				Converter:  mock_converter.NewMockConverter(gomock.NewController(t)),
			},
			args: args{
				ctx:    context.Background(),
				phrase: testPhrase,
			},
			wantErr: true,
			prepare: func(f *fields) {
				f.Converter.EXPECT().Convert(testContent, "wav").Return(nil, errors.New("failed to convert"))
			},
		},
		{
			name: "failed to save converted audio files",
			fields: fields{
				Repository: mock_repository.NewMockPhraseRepository(gomock.NewController(t)),
				Storage:    mock_storage.NewMockAudioStorage(gomock.NewController(t)),
				Converter:  mock_converter.NewMockConverter(gomock.NewController(t)),
			},
			args: args{
				ctx:    context.Background(),
				phrase: testPhrase,
			},
			wantErr: true,
			prepare: func(f *fields) {
				f.Converter.EXPECT().Convert(testContent, "wav").Return(testContent, nil)
				f.Storage.EXPECT().SaveFile(context.Background(), defaultUserIDinStr, gomock.Any(), testPhrase).
					Return("", errors.New("some os error"))
			},
		},
		{
			name: "failed to save metadata to DB, file is deleted",
			fields: fields{
				Repository: mock_repository.NewMockPhraseRepository(gomock.NewController(t)),
				Storage:    mock_storage.NewMockAudioStorage(gomock.NewController(t)),
				Converter:  mock_converter.NewMockConverter(gomock.NewController(t)),
			},
			args: args{
				ctx:    context.Background(),
				phrase: testPhrase,
			},
			wantErr: true,
			prepare: func(f *fields) {
				f.Converter.EXPECT().Convert(testContent, "wav").Return(testContent, nil)
				f.Storage.EXPECT().SaveFile(context.Background(), defaultUserIDinStr, gomock.Any(), testPhrase).Return(savedPath, nil)
				f.Repository.EXPECT().SavePhraseMetadata(context.Background(), testPhraseWithPath).Return(errors.New("db error"))
				f.Storage.EXPECT().DeleteFile(context.Background(), savedPath).Return(nil)
			},
		},
		{
			name: "failed to save metadata to DB, file is not deleted",
			fields: fields{
				Repository: mock_repository.NewMockPhraseRepository(gomock.NewController(t)),
				Storage:    mock_storage.NewMockAudioStorage(gomock.NewController(t)),
				Converter:  mock_converter.NewMockConverter(gomock.NewController(t)),
			},
			args: args{
				ctx:    context.Background(),
				phrase: testPhrase,
			},
			wantErr: true,
			prepare: func(f *fields) {
				f.Converter.EXPECT().Convert(testContent, "wav").Return(testContent, nil)
				f.Storage.EXPECT().SaveFile(context.Background(), defaultUserIDinStr, gomock.Any(), testPhrase).Return(savedPath, nil)
				f.Repository.EXPECT().SavePhraseMetadata(context.Background(), testPhraseWithPath).Return(errors.New("db error"))
				f.Storage.EXPECT().DeleteFile(context.Background(), savedPath).Return(errors.New("some error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &PhraseServiceProvider{
				Repository: tt.fields.Repository,
				Storage:    tt.fields.Storage,
				Converter:  tt.fields.Converter,
			}
			tt.prepare(&tt.fields)
			if err := ps.Store(tt.args.ctx, tt.args.phrase); (err != nil) != tt.wantErr {
				t.Errorf("PhraseServiceProvider.Store() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPhraseServiceProvider_Retrieve(t *testing.T) {
	type fields struct {
		Repository *mock_repository.MockPhraseRepository
		Storage    *mock_storage.MockAudioStorage
		Converter  *mock_converter.MockConverter
	}
	type args struct {
		ctx       context.Context
		userID    int
		phraseID  int
		extension string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Phrase
		wantErr bool
		prepare func(f *fields)
	}{
		{
			name: "happy path",
			fields: fields{
				Repository: mock_repository.NewMockPhraseRepository(gomock.NewController(t)),
				Storage:    mock_storage.NewMockAudioStorage(gomock.NewController(t)),
				Converter:  mock_converter.NewMockConverter(gomock.NewController(t)),
			},
			args: args{
				ctx:       context.Background(),
				userID:    defaultUserID,
				phraseID:  defaultPhraseID,
				extension: ext,
			},
			want: model.Phrase{
				ID:       defaultPhraseID,
				UserID:   defaultUserID,
				Content:  testContent,
				FilePath: savedPath,
			},
			wantErr: false,
			prepare: func(f *fields) {
				f.Repository.EXPECT().FindPhraseMetadata(context.Background(), defaultUserID, defaultPhraseID).Return(model.Phrase{
					ID:       defaultPhraseID,
					UserID:   defaultUserID,
					FilePath: savedPath,
				}, nil)
				f.Storage.EXPECT().LoadFile(context.Background(), savedPath).Return(testContent, nil)
				f.Converter.EXPECT().Convert(testContent, ext).Return(testContent, nil)
			},
		},
		{
			name: "error converting saved file",
			fields: fields{
				Repository: mock_repository.NewMockPhraseRepository(gomock.NewController(t)),
				Storage:    mock_storage.NewMockAudioStorage(gomock.NewController(t)),
				Converter:  mock_converter.NewMockConverter(gomock.NewController(t)),
			},
			args: args{
				ctx:       context.Background(),
				userID:    defaultUserID,
				phraseID:  defaultPhraseID,
				extension: ext,
			},
			want:    model.Phrase{},
			wantErr: true,
			prepare: func(f *fields) {
				f.Repository.EXPECT().FindPhraseMetadata(context.Background(), defaultUserID, defaultPhraseID).Return(model.Phrase{
					ID:       defaultPhraseID,
					UserID:   defaultUserID,
					FilePath: savedPath,
				}, nil)
				f.Storage.EXPECT().LoadFile(context.Background(), savedPath).Return(testContent, nil)
				f.Converter.EXPECT().Convert(testContent, ext).Return(nil, errors.New("converting error"))
			},
		},
		{
			name: "error loading saved file",
			fields: fields{
				Repository: mock_repository.NewMockPhraseRepository(gomock.NewController(t)),
				Storage:    mock_storage.NewMockAudioStorage(gomock.NewController(t)),
				Converter:  mock_converter.NewMockConverter(gomock.NewController(t)),
			},
			args: args{
				ctx:       context.Background(),
				userID:    defaultUserID,
				phraseID:  defaultPhraseID,
				extension: ext,
			},
			want:    model.Phrase{},
			wantErr: true,
			prepare: func(f *fields) {
				f.Repository.EXPECT().FindPhraseMetadata(context.Background(), defaultUserID, defaultPhraseID).Return(model.Phrase{
					ID:       defaultPhraseID,
					UserID:   defaultUserID,
					FilePath: savedPath,
				}, nil)
				f.Storage.EXPECT().LoadFile(context.Background(), savedPath).Return(testContent, errors.New("file corrupted"))
			},
		},
		{
			name: "error invalid file path",
			fields: fields{
				Repository: mock_repository.NewMockPhraseRepository(gomock.NewController(t)),
				Storage:    mock_storage.NewMockAudioStorage(gomock.NewController(t)),
				Converter:  mock_converter.NewMockConverter(gomock.NewController(t)),
			},
			args: args{
				ctx:       context.Background(),
				userID:    defaultUserID,
				phraseID:  defaultPhraseID,
				extension: ext,
			},
			want:    model.Phrase{},
			wantErr: true,
			prepare: func(f *fields) {
				f.Repository.EXPECT().FindPhraseMetadata(context.Background(), defaultUserID, defaultPhraseID).Return(model.Phrase{
					ID:       defaultPhraseID,
					UserID:   defaultUserID,
					FilePath: "",
				}, nil)
			},
		},
		{
			name: "error loading phrase metadata",
			fields: fields{
				Repository: mock_repository.NewMockPhraseRepository(gomock.NewController(t)),
				Storage:    mock_storage.NewMockAudioStorage(gomock.NewController(t)),
				Converter:  mock_converter.NewMockConverter(gomock.NewController(t)),
			},
			args: args{
				ctx:       context.Background(),
				userID:    defaultUserID,
				phraseID:  defaultPhraseID,
				extension: ext,
			},
			want:    model.Phrase{},
			wantErr: true,
			prepare: func(f *fields) {
				f.Repository.EXPECT().FindPhraseMetadata(context.Background(), defaultUserID, defaultPhraseID).Return(model.Phrase{}, errors.New("db error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &PhraseServiceProvider{
				Repository: tt.fields.Repository,
				Storage:    tt.fields.Storage,
				Converter:  tt.fields.Converter,
			}
			tt.prepare(&tt.fields)
			got, err := ps.Retrieve(tt.args.ctx, tt.args.userID, tt.args.phraseID, tt.args.extension)
			if (err != nil) != tt.wantErr {
				t.Errorf("PhraseServiceProvider.Retrieve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PhraseServiceProvider.Retrieve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPhraseServiceProvider(t *testing.T) {
	type args struct {
		phraseRepo     *mock_repository.MockPhraseRepository
		audioConverter *mock_converter.MockConverter
		storage        *mock_storage.MockAudioStorage
	}
	tests := []struct {
		name string
		args args
		want *PhraseServiceProvider
	}{
		{
			name: "constructor",
			args: args{
				phraseRepo:     mock_repository.NewMockPhraseRepository(gomock.NewController(t)),
				storage:        mock_storage.NewMockAudioStorage(gomock.NewController(t)),
				audioConverter: mock_converter.NewMockConverter(gomock.NewController(t)),
			},
			want: &PhraseServiceProvider{
				Repository: mock_repository.NewMockPhraseRepository(gomock.NewController(t)),
				Storage:    mock_storage.NewMockAudioStorage(gomock.NewController(t)),
				Converter:  mock_converter.NewMockConverter(gomock.NewController(t)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPhraseServiceProvider(tt.args.phraseRepo, tt.args.audioConverter, tt.args.storage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPhraseServiceProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func prepareTestAudioContent() []byte {
	f, err := os.Open("../../../fixtures/audio-test.mp3")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer f.Close()

	c, _ := io.ReadAll(f)
	return c
}
