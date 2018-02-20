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

func getKeukendienst(c echo.Context) error {
	keukendienst := NewKeukendienst()

	err := getSetting("keukendienst", &keukendienst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, keukendienst)
}

func setKeukendienst(c echo.Context) error {
	keukendienst := NewKeukendienst()
	c.Bind(&keukendienst)
	err := editSettings("keukendienst", keukendienst)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	go sendUpdate()
	return c.JSON(http.StatusOK, keukendienst)
}

func setTimers() {
	resetTimers()
	news, _ := getNewsItems()
	for _, item := range news {
		addTimer(item.From)
		addTimer(item.To)
	}
}
