package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "data", nil)
	})
	e.GET("/healthcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	return e
}

func SetupRoute(e *echo.Echo, h HTTPHandler) {
	r := e.Group("/audio")
	// r.Use(m)

	r.POST("/user/:user_id/phrase/:phrase_id", h.SavePhrase)
	r.GET("/user/:user_id/phrase/:phrase_id/:extension", h.RetrievePhrase)
}
