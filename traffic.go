package main

import (
	"net/http"

	"github.com/SlyMarbo/rss"
	"github.com/labstack/echo"
)

func getTrafficHandler(c echo.Context) error {
	feed, err := rss.Fetch("http://www.verkeerscentrum.be/rss/100-INC%7CLOS%7CINF%7CPEVT.xml")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, feed.Items)
}
