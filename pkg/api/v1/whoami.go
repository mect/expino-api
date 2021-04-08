package v1

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/mect/expino-api/pkg/api/auth"

	"github.com/labstack/echo/v4"
)

func (h *HTTPHandler) checkAuth(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.Claim)
	if claims.Name == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"status": "JWT incorrect"})
	}
	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}
