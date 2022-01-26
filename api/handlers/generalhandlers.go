package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "hello from the web side!")
}
