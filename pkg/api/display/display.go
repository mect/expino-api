package display

import (
	"net/http"

	"github.com/mect/expino-api/pkg/db"

	"github.com/labstack/echo/v4"
)

func (h *HTTPHandler) handleDisplay(c echo.Context) error {
	display, ok := c.Get("display").(*db.Display)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "error reading display data from session"})
	}

	display.Token = "" // would be silly to return this

	return c.JSON(http.StatusOK, display)
}
