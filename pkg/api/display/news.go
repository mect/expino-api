package display

import (
	"fmt"
	"net/http"
	"time"

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

	var displayItems []db.NewsItem

	for _, item := range newsItems {
		if item.From != nil && item.To != nil {
			if time.Now().Before(*item.From) {
				continue
			}
			if time.Now().After(*item.To) {
				continue
			}
		}

		displayItems = append(displayItems, item)
	}

	return c.JSON(http.StatusOK, displayItems)
}
