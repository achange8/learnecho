package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

func MainJWT(c echo.Context) error {

	return c.JSON(http.StatusOK, " you are on the secrt page in jwt")
}
