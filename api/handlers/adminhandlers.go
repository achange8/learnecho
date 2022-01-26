package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

func MainAdmin(g echo.Context) error {
	return g.String(http.StatusOK, "horay you are on secret admin main page!")
}
