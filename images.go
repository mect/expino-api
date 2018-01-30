package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
	filetype "gopkg.in/h2non/filetype.v1"
)

func uploadImage(c echo.Context) error {
	hostname := c.FormValue("host")

	// Source
	file, err := c.FormFile("image")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	b, err := ioutil.ReadAll(src)

	sum := fmt.Sprintf("%x", sha256.Sum256(b))

	saveFile(sum, b)

	return c.JSON(http.StatusOK, map[string]string{"link": hostname + "/api/images/" + sum})
}

func getImage(c echo.Context) error {
	b, err := getFile(c.Param("name"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	contentType := ""
	kind, unknown := filetype.Match(b)
	if unknown != nil {
		contentType = "image/jpeg" // best bet
	} else {
		contentType = kind.MIME.Value
	}

	return c.Blob(http.StatusOK, contentType, b)
}
