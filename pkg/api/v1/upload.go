package v1

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/labstack/echo/v4"
)

const STATICDIR = "expino-static"

func (h *HTTPHandler) handleUpload(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("opening file: %v", err)})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("opening file: %v", err)})
	}
	defer src.Close()

	filenameParts := strings.Split(file.Filename, ".")
	if len(filenameParts) < 2 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "File needs a file extension"})
	}

	tmpFile := path.Join(STATICDIR, file.Filename)
	dst, err := os.Create(tmpFile)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("opening creating tmp file: %v", err)})
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("opening writing tmp file: %v", err)})
	}

	hash, err := getHash(tmpFile)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("opening creating hash of file: %v", err)})
	}

	newName := fmt.Sprintf("%s.%s", hash, filenameParts[len(filenameParts)-1])
	err = os.Rename(tmpFile, path.Join(STATICDIR, newName))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("error moving file: %v", err)})
	}

	return c.JSON(http.StatusOK, echo.Map{"location": fmt.Sprintf("https://%s/static/%s", c.Request().Host, newName)})
}

func getHash(file string) (string, error) {
	hasher := sha256.New()

	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
