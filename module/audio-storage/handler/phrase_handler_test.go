package handler

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dibikhairurrazi/audio-storage/module/audio-storage/model"
	mock_service "github.com/dibikhairurrazi/audio-storage/module/audio-storage/test/mock/module/audio-storage/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/mock/gomock"
)

func TestHTTPHandler_SavePhrase(t *testing.T) {
	type fields struct {
		UserService  *mock_service.MockUserService
		AudioService *mock_service.MockPhraseService
	}
	type args struct {
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		prepare func(f *fields)
		setup   func() echo.Context
	}{
		{
			name: "happy path",
			fields: fields{
				UserService:  mock_service.NewMockUserService(gomock.NewController(t)),
				AudioService: mock_service.NewMockPhraseService(gomock.NewController(t)),
			},
			args:    args{},
			wantErr: false,
			prepare: func(f *fields) {
				f.UserService.EXPECT().FindUser(gomock.Any(), 1).Return(model.User{ID: 1}, nil)
				f.AudioService.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func() echo.Context {
				path := "../../../fixtures/audio-test.mp3"
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("audio_file", path)
				sample, _ := os.Open(path)

				io.Copy(part, sample)
				writer.Close()
				sample.Close()

				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/user/:user_id/phrase/:phrase_id", body)
				req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetParamNames("user_id", "phrase_id")
				c.SetParamValues("1", "1")

				return c
			},
		},
		{
			name: "invalid user ID param",
			fields: fields{
				UserService:  mock_service.NewMockUserService(gomock.NewController(t)),
				AudioService: mock_service.NewMockPhraseService(gomock.NewController(t)),
			},
			args:    args{},
			wantErr: true,
			prepare: func(f *fields) {

			},
			setup: func() echo.Context {
				path := "../../../fixtures/audio-test.mp3"
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("audio_file", path)
				sample, _ := os.Open(path)

				io.Copy(part, sample)
				writer.Close()
				sample.Close()

				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/user/:user_id/phrase/:phrase_id", body)
				req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetParamNames("user_id", "phrase_id")
				c.SetParamValues("abc", "1")

				return c
			},
		},
		{
			name: "invalid phrase ID param",
			fields: fields{
				UserService:  mock_service.NewMockUserService(gomock.NewController(t)),
				AudioService: mock_service.NewMockPhraseService(gomock.NewController(t)),
			},
			args:    args{},
			wantErr: true,
			prepare: func(f *fields) {

			},
			setup: func() echo.Context {
				path := "../../../fixtures/audio-test.mp3"
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("audio_file", path)
				sample, _ := os.Open(path)

				io.Copy(part, sample)
				writer.Close()
				sample.Close()

				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/user/:user_id/phrase/:phrase_id", body)
				req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetParamNames("user_id", "phrase_id")
				c.SetParamValues("1", "abc")

				return c
			},
		},
		{
			name: "no audio file",
			fields: fields{
				UserService:  mock_service.NewMockUserService(gomock.NewController(t)),
				AudioService: mock_service.NewMockPhraseService(gomock.NewController(t)),
			},
			args:    args{},
			wantErr: true,
			prepare: func(f *fields) {
				f.UserService.EXPECT().FindUser(gomock.Any(), 1).Return(model.User{ID: 1}, nil)
			},
			setup: func() echo.Context {
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/user/:user_id/phrase/:phrase_id", nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetParamNames("user_id", "phrase_id")
				c.SetParamValues("1", "1")

				return c
			},
		},
		{
			name: "user not found",
			fields: fields{
				UserService:  mock_service.NewMockUserService(gomock.NewController(t)),
				AudioService: mock_service.NewMockPhraseService(gomock.NewController(t)),
			},
			args:    args{},
			wantErr: true,
			prepare: func(f *fields) {
				f.UserService.EXPECT().FindUser(gomock.Any(), 1).Return(model.User{}, errors.New("not found"))
			},
			setup: func() echo.Context {
				path := "../../../fixtures/audio-test.mp3"
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("audio_file", path)
				sample, _ := os.Open(path)

				io.Copy(part, sample)
				writer.Close()
				sample.Close()

				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/user/:user_id/phrase/:phrase_id", body)
				req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetParamNames("user_id", "phrase_id")
				c.SetParamValues("1", "1")

				return c
			},
		},
		{
			name: "error storing audio",
			fields: fields{
				UserService:  mock_service.NewMockUserService(gomock.NewController(t)),
				AudioService: mock_service.NewMockPhraseService(gomock.NewController(t)),
			},
			args:    args{},
			wantErr: true,
			prepare: func(f *fields) {
				f.UserService.EXPECT().FindUser(gomock.Any(), 1).Return(model.User{ID: 1}, nil)
				f.AudioService.EXPECT().Store(gomock.Any(), gomock.Any()).Return(errors.New("error"))
			},
			setup: func() echo.Context {
				path := "../../../fixtures/audio-test.mp3"
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("audio_file", path)
				sample, _ := os.Open(path)

				io.Copy(part, sample)
				writer.Close()
				sample.Close()

				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/user/:user_id/phrase/:phrase_id", body)
				req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetParamNames("user_id", "phrase_id")
				c.SetParamValues("1", "1")

				return c
			},
		},
		{
			name: "not an audio file",
			fields: fields{
				UserService:  mock_service.NewMockUserService(gomock.NewController(t)),
				AudioService: mock_service.NewMockPhraseService(gomock.NewController(t)),
			},
			args:    args{},
			wantErr: true,
			prepare: func(f *fields) {
				f.UserService.EXPECT().FindUser(gomock.Any(), 1).Return(model.User{ID: 1}, nil)
			},
			setup: func() echo.Context {
				path := "../../../fixtures/pdf-test.mp3"
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("audio_file", path)
				sample, _ := os.Open(path)

				io.Copy(part, sample)
				writer.Close()
				sample.Close()

				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/user/:user_id/phrase/:phrase_id", body)
				req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetParamNames("user_id", "phrase_id")
				c.SetParamValues("1", "1")

				return c
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HTTPHandler{
				UserService:  tt.fields.UserService,
				AudioService: tt.fields.AudioService,
			}
			c := tt.setup()
			tt.prepare(&tt.fields)
			if err := h.SavePhrase(c); (err != nil) != tt.wantErr {
				t.Errorf("HTTPHandler.SavePhrase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHTTPHandler_RetrievePhrase(t *testing.T) {
	type fields struct {
		UserService  *mock_service.MockUserService
		AudioService *mock_service.MockPhraseService
	}
	type args struct {
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		prepare func(f *fields)
		setup   func() echo.Context
	}{
		{
			name: "happy path",
			fields: fields{
				UserService:  mock_service.NewMockUserService(gomock.NewController(t)),
				AudioService: mock_service.NewMockPhraseService(gomock.NewController(t)),
			},
			args:    args{},
			wantErr: false,
			prepare: func(f *fields) {
				f.UserService.EXPECT().FindUser(gomock.Any(), 1).Return(model.User{ID: 1}, nil)
				f.AudioService.EXPECT().Retrieve(gomock.Any(), 1, 1, "mp3").Return(model.Phrase{
					ID:       1,
					UserID:   1,
					FilePath: "/1/audio-test.mp3",
				}, nil)
			},
			setup: func() echo.Context {
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/user/:user_id/phrase/:phrase_id/:extension", nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetParamNames("user_id", "phrase_id", "extension")
				c.SetParamValues("1", "1", "mp3")

				return c
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HTTPHandler{
				UserService:  tt.fields.UserService,
				AudioService: tt.fields.AudioService,
			}
			c := tt.setup()
			tt.prepare(&tt.fields)
			if err := h.RetrievePhrase(c); (err != nil) != tt.wantErr {
				t.Errorf("HTTPHandler.RetrievePhrase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
