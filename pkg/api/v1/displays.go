package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"gorm.io/gorm/clause"

	"github.com/labstack/echo/v4"
	"github.com/mect/expino-api/pkg/db"
)

func (h *HTTPHandler) handleDisplayList(c echo.Context) error {
	var displays []db.Display
	res := h.db.Preload(clause.Associations).Order("id").Find(&displays)
	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("error reading data: %v", res.Error)})
	}

	return c.JSON(http.StatusOK, displays)
}

func (h *HTTPHandler) handleDisplayGet(c echo.Context) error {
	display := &db.Display{}
	res := h.db.Preload(clause.Associations).First(display, "id = ?", c.Param("id"))
	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("error reading data: %v", res.Error)})
	}

	return c.JSON(http.StatusOK, display)
}

func (h *HTTPHandler) handleDisplayCreate(c echo.Context) error {
	display := &db.Display{}
	err := c.Bind(display)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("error reading data: %v", err)})
	}

	display.ID = 0

	res := h.db.Create(display)
	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("error saving data: %v", res.Error)})
	}

	return c.JSON(http.StatusOK, display)
}

func (h *HTTPHandler) handleDisplayUpdate(c echo.Context) error {
	display := &db.Display{}
	err := c.Bind(display)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("error reading data: %v", err)})
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("error reading ID: %v", err)})
	} else if id < 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "ID cannot be negative"})
	}
	display.ID = uint(id)

	res := h.db.Save(display)
	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("error saving data: %v", res.Error)})
	}

	h.broadcastUpdate()

	return c.JSON(http.StatusOK, display)
}

func (h *HTTPHandler) handleDisplayDelete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("error reading ID: %v", err)})
	} else if id < 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "ID cannot be negative"})
	}

	res := h.db.Delete(&db.Display{
		Model: gorm.Model{
			ID: uint(id),
		},
	})
	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("error deleting data: %v", res.Error)})
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}
