package display

import (
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/mect/expino-api/pkg/db"

	"github.com/labstack/echo/v4"
)

// wish I could prevent this but we need proper CORS on this one
const delijnHost = "https://b2cservices.delijn.be"

func (h *HTTPHandler) handleDelijn(c echo.Context) error {
	display, ok := c.Get("display").(*db.Display)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "error reading display data from session"})
	}

	if display.TickerRSS == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "no RSS set"})
	}

	resp, err := http.Get(fmt.Sprintln("%s%s", delijnHost, path.Join("/", c.QueryParam("path"))))
	if err != nil {
		return err
	}

	c.Response().WriteHeader(resp.StatusCode)
	io.Copy(c.Response(), resp.Body)

	return nil
}
