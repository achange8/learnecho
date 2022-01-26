package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

func MainCookie(g echo.Context) error {
	return g.String(http.StatusOK, "you are on secret cookie main page!")
}
