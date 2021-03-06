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
	display, ok := c.Get("display").(*db.Display)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "error reading display data from session"})
	}

	var newsItems []db.NewsItem
	res := h.db.Preload(clause.Associations).Where("display_id", display.ID).Where("hidden", false).Order("\"order\",id").Find(&newsItems)
	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("error reading data: %v", res.Error)})
	}

	var displayItems []db.NewsItem

	for _, item := range newsItems {
		showNow := len(item.TimeFrames) == 0
		for _, tf := range item.TimeFrames {
			if time.Now().After(tf.From) && time.Now().Before(tf.To) {
				showNow = true
				break
			}
		}

		if showNow {
			displayItems = append(displayItems, item)
		}
	}

	return c.JSON(http.StatusOK, displayItems)
}
