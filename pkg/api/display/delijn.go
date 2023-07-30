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
	req, err := http.NewRequest("GET", delijnHost+path.Join("/", c.QueryParam("path")), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", c.Request().Header.Get("Ocp-Apim-Subscription-Key"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	c.Response().WriteHeader(resp.StatusCode)
	io.Copy(c.Response(), resp.Body)

	return nil
}
