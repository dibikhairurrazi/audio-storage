package handler

import (
	"database/sql"
	"errors"
	"io"
	"net/http"
	"slices"
	"strconv"

	"github.com/dibikhairurrazi/audio-storage/module/audio-storage/model"
	"github.com/dibikhairurrazi/audio-storage/module/audio-storage/service"
	"github.com/labstack/echo/v4"
)

var (
	allowedExtension = []string{"mp3"} // to allow for more extension, insert into this array
)

type HTTPHandler struct {
	UserService  service.UserService
	AudioService service.PhraseService
}

func NewHTTPHandler(ps service.PhraseService, us service.UserService) HTTPHandler {
	return HTTPHandler{
		AudioService: ps,
		UserService:  us,
	}
}

// SavePhrase handle save phrase endpoint, it will check if the user valid or not first and then upsert the phrase.
// assuming that if the user use the same phrase ID, he/she want to overwrite it
func (h *HTTPHandler) SavePhrase(c echo.Context) error {
	// TODO authorization and validate user_id with auth token

	userID, phraseID, err := validateAndRetrieveUserIDandPhraseID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid parameter"))
	}

	user, err := h.UserService.FindUser(c.Request().Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("user not found"))
	}

	file, err := c.FormFile("audio_file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	f, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	defer func() {
		_ = f.Close()
	}()

	content, err := io.ReadAll(f)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if !isValidFileType(content) {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("not an audio file"))
	}

	phrase := model.Phrase{
		ID:               phraseID,
		UserID:           user.ID,
		Content:          content,
		OriginalFileName: file.Filename,
	}

	err = h.AudioService.Store(c.Request().Context(), phrase)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, model.CommonResponse{
		Data: "phrase saved succesfully",
	})
}

// RetrievePhrase handle retrieving endpoint, it need 3 param as its path, user ID, phrase ID, and extension
func (h *HTTPHandler) RetrievePhrase(c echo.Context) error {
	// TODO authorization and validate user_id with auth token

	userID, phraseID, err := validateAndRetrieveUserIDandPhraseID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid parameter"))
	}

	ext := c.Param("extension")

	if !isValidExtension(ext) {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid extension, please use mp3"))
	}

	user, err := h.UserService.FindUser(c.Request().Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("user not found"))
	}

	audio, err := h.AudioService.Retrieve(c.Request().Context(), user.ID, phraseID, ext)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.Attachment(audio.FilePath, audio.OriginalFileName)
}

func validateAndRetrieveUserIDandPhraseID(c echo.Context) (int, int, error) {
	userIDParam := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		return 0, 0, err
	}

	phraseIDParam := c.Param("phrase_id")
	phraseID, err := strconv.ParseInt(phraseIDParam, 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return int(userID), int(phraseID), nil
}

func isValidFileType(file []byte) bool {
	fileType := http.DetectContentType(file)
	// return strings.HasPrefix("audio/") // use this return for enabling other audio files
	return fileType == "audio/mpeg" // Only allow audio mp3 files
}

func isValidExtension(s string) bool {
	return slices.Contains(allowedExtension, s)
}
