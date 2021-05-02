package display

import (
	"io"
	"net/http"

	"github.com/mect/expino-api/pkg/db"

	"github.com/labstack/echo/v4"
)

func (h *HTTPHandler) handleRss(c echo.Context) error {
	display, ok := c.Get("display").(*db.Display)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "error reading display data from session"})
	}

	if display.TickerRSS == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "no RSS set"})
	}

	resp, err := http.Get(display.TickerRSS)
	if err != nil {
		return err
	}

	c.Response().WriteHeader(http.StatusOK)
	io.Copy(c.Response(), resp.Body)

	return nil
}
