package display

import (
	"github.com/mect/expino-api/pkg/db"

	"github.com/labstack/echo/v4"
)

type HTTPHandler struct {
	db *db.Connection
}

func NewHTTPHandler(db *db.Connection) *HTTPHandler {
	return &HTTPHandler{
		db: db,
	}
}

func (h *HTTPHandler) Register(e *echo.Echo) {
	e.Use(h.checkAuth)

	e.GET("/display", h.handleDisplay)
	e.GET("/display/news", h.handleNewsList)
	e.GET("/display/rss", h.handleRss)
}
