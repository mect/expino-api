package display

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

const RSSurl = "https://www.vrt.be/vrtnws/nl.rss.articles.xml"

func (h *HTTPHandler) handleRss(c echo.Context) error {
	resp, err := http.Get(RSSurl)
	if err != nil {
		return err
	}

	c.Response().WriteHeader(http.StatusOK)
	io.Copy(c.Response(), resp.Body)

	return nil
}
