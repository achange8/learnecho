package handlers

import (
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
)

func MainJWT(c echo.Context) error {
	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	log.Println("User name:", claims["name"], "User ID:", claims["jti"])

	return c.JSON(http.StatusOK, " you are on the secrt page in jwt")
}
