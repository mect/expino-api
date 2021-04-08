package display

import (
	"fmt"
	"net/http"
	"strings"

	"gorm.io/gorm"

	"github.com/mect/expino-api/pkg/db"

	"github.com/labstack/echo/v4"
)

// checkAuth checks for the display token in the API
func (h *HTTPHandler) checkAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !strings.HasPrefix(c.Path(), "/display") {
			return next(c)
		}
		auth := c.Request().Header.Get("Authorization")
		if auth == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "authorization header not set"})
		}

		display := &db.Display{}
		res := h.db.First(display, "token = ?", strings.Replace(auth, "Bearer ", "", -1))

		if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("error looking up token: %v", res.Error)})
		}

		if res.Error == gorm.ErrRecordNotFound || display.Token == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "authorization header incorrect"})
		}

		c.Set("display", display)

		return next(c)
	}
}
