package v1

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/mect/expino-api/pkg/db"

	"github.com/labstack/echo/v4"
)

type HTTPHandler struct {
	db *db.Connection
	io *socketio.Server
}

func NewHTTPHandler(db *db.Connection, io *socketio.Server) *HTTPHandler {
	return &HTTPHandler{
		db: db,
		io: io,
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

	e.POST("/v1/static/upload", h.handleUpload)
}

func (h *HTTPHandler) broadcastUpdate() {
	h.io.BroadcastToRoom("/", "updates", "update")
}
