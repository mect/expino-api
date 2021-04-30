package v1

import (
	"github.com/mect/expino-api/pkg/db"

	"github.com/labstack/echo/v4"
)

type HTTPHandler struct {
	db        *db.Connection
	broadcast func(data string)
}

func NewHTTPHandler(db *db.Connection, broadcast func(data string)) *HTTPHandler {
	return &HTTPHandler{
		db:        db,
		broadcast: broadcast,
	}
}

func (h *HTTPHandler) Register(e *echo.Echo) {
	// whoami
	e.GET("/v1/auth/check", h.checkAuth)

	// news
	e.GET("/v1/news/items", h.handleNewsList)
	e.GET("/v1/news/item/:id", h.handleNewsGet)

	e.POST("/v1/news/item", h.handleNewsCreate)
	e.POST("/v1/news/item/:id", h.handleNewsUpdate)
	e.DELETE("/v1/news/item/:id", h.handleNewsDelete)

	// displays
	e.GET("/v1/display/items", h.handleDisplayList)
	e.GET("/v1/display/item/:id", h.handleDisplayGet)

	e.POST("/v1/display/item", h.handleDisplayCreate)
	e.POST("/v1/display/item/:id", h.handleDisplayUpdate)
	e.DELETE("/v1/display/item/:id", h.handleDisplayDelete)

	// static
	e.POST("/v1/static/upload", h.handleUpload)
}

func (h *HTTPHandler) broadcastUpdate() {
	h.broadcast("UPDATE")
}
