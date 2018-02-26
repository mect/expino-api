package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func getNewsHandler(c echo.Context) error {
	i, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	news, err := getNewsItem(int(i))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, news)
}

func getAllNewsHandler(c echo.Context) error {
	news, err := getNewsItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, news)
}

func getCurrentNewsHandler(c echo.Context) error {
	news, err := getNewsItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	now := time.Now()
	currentNews := []NewsItem{}
	for _, item := range news {
		if now.After(item.From.Truncate(time.Second)) && now.Before(item.To.Truncate(time.Second)) {
			currentNews = append(currentNews, item)
		}
	}

	return c.JSON(http.StatusOK, currentNews)
}

func addNewsHandler(c echo.Context) error {
	item := NewNewsItem()
	c.Bind(&item)
	err := addNewsItem(item)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	go sendUpdate()
	return c.JSON(http.StatusOK, item)
}

func editNewsHandler(c echo.Context) error {
	item := NewNewsItem()
	c.Bind(&item)
	err := editNewsItem(item)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	go sendUpdate()
	return c.JSON(http.StatusOK, item)
}

func deleteNewsHandler(c echo.Context) error {
	i, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	err = deleteNewsItem(int(i))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	go sendUpdate()
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func editFeatureSlides(c echo.Context) error {
	items := []string{}
	c.Bind(&items)
	err := editSettings("featureSlides", items)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	setTimers()
	go sendUpdate()
	return c.JSON(http.StatusOK, items)
}

func getFeatureSlides(c echo.Context) error {
	items := []string{}
	err := getSetting("featureSlides", &items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, items)
}

func getTickerItemsHandler(c echo.Context) error {
	items, err := getTickerItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, items)
}

func addTickerItemHandler(c echo.Context) error {
	item := TickerItem{}
	c.Bind(&item)
	_, err := time.ParseDuration(item.Back)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	_, err = time.ParseDuration(item.Interval)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err = addTickerItem(item)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	go sendUpdate()
	return c.JSON(http.StatusOK, item)
}

func deleteTickerItemsHandler(c echo.Context) error {
	i, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	err = deleteTickerItem(int(i))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	go sendUpdate()
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func getGraphItemsHandler(c echo.Context) error {
	items, err := getGraphItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, items)
}

func addGraphItemHandler(c echo.Context) error {
	item := GraphItem{}
	c.Bind(&item)
	err := addGraphItem(item)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	go sendUpdate()
	return c.JSON(http.StatusOK, item)
}

func deleteGraphItemsHandler(c echo.Context) error {
	i, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	err = deleteGraphItem(int(i))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	go sendUpdate()
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func getKeukenDienstItemsHandler(c echo.Context) error {
	items, err := getKeukenDienstItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, items)
}

func addKeukenDienstItemsHandler(c echo.Context) error {
	item := KeukendienstItem{}
	c.Bind(&item)
	err := addKeukenDienstItems(item)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, item)
}

func deleteKeukenDienstItemsHandler(c echo.Context) error {
	i, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	err = deleteKeukenDienstItems(int(i))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func getCurrentKeukenDienstItemHandler(c echo.Context) error {
	news, err := getKeukenDienstItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	now := time.Now()
	current := KeukendienstItem{}
	for _, item := range news {
		if now.After(item.From.Truncate(time.Second)) && now.Before(item.To.Truncate(time.Second)) {
			current = item
		}
	}

	return c.JSON(http.StatusOK, current)
}

func setTimers() {
	resetTimers()
	news, _ := getNewsItems()
	for _, item := range news {
		addTimer(item.From)
		addTimer(item.To)
	}
}
