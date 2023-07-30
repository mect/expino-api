package display

import (
	"io"
	"net/http"
	"path"

	"github.com/labstack/echo/v4"
)

// wish I could prevent this but we need proper CORS on this one
const delijnHost = "https://api.delijn.be"

func (h *HTTPHandler) handleDelijn(c echo.Context) error {
	resp, err := http.Get(delijnHost + path.Join("/", c.QueryParam("path")))
	if err != nil {
		return err
	}

	c.Response().WriteHeader(resp.StatusCode)
	io.Copy(c.Response(), resp.Body)

	return nil
}
