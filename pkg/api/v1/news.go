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

func (h *HTTPHandler) handleNewsList(c echo.Context) error {
	var newsItems []db.NewsItem
	res := h.db.Preload(clause.Associations).Find(&newsItems)
	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("error reading data: %v", res.Error)})
	}

	return c.JSON(http.StatusOK, newsItems)
}

func (h *HTTPHandler) handleNewsGet(c echo.Context) error {
	newsItem := &db.NewsItem{}
	res := h.db.Preload(clause.Associations).First(newsItem, "id = ?", c.Param("id"))
	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("error reading data: %v", res.Error)})
	}

	return c.JSON(http.StatusOK, newsItem)
}

func (h *HTTPHandler) handleNewsCreate(c echo.Context) error {
	newsItem := &db.NewsItem{}
	err := c.Bind(newsItem)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("error reading data: %v", err)})
	}

	newsItem.ID = 0

	res := h.db.Create(newsItem)
	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("error saving data: %v", res.Error)})
	}

	h.broadcastUpdate()

	return c.JSON(http.StatusOK, newsItem)
}

func (h *HTTPHandler) handleNewsUpdate(c echo.Context) error {
	newsItem := &db.NewsItem{}
	err := c.Bind(newsItem)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("error reading data: %v", err)})
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("error reading ID: %v", err)})
	} else if id < 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "ID cannot be negative"})
	}
	newsItem.ID = uint(id)

	res := h.db.Save(newsItem)
	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("error saving data: %v", res.Error)})
	}

	err = h.db.Model(&newsItem).Association("LanguageItems").Replace(&newsItem.LanguageItems)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("error saving data: %v", res.Error)})

	}

	for _, langItem := range newsItem.LanguageItems {
		res := h.db.Save(&langItem)
		if res.Error != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("error saving data: %v", res.Error)})
		}
	}

	h.broadcastUpdate()

	return c.JSON(http.StatusOK, newsItem)
}

func (h *HTTPHandler) handleNewsDelete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("error reading ID: %v", err)})
	} else if id < 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "ID cannot be negative"})
	}

	res := h.db.Delete(&db.NewsItem{
		Model: gorm.Model{
			ID: uint(id),
		},
	})
	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("error deleting data: %v", res.Error)})
	}

	h.broadcastUpdate()

	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}
