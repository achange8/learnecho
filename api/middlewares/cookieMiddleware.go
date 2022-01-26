package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

func SetMCookieMiddleware(g *echo.Group) {

	g.Use(checkCookie)

}
func checkCookie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("sessionID")
		if err != nil {
			if strings.Contains(err.Error(), "named cookie not present") {
				return c.String(http.StatusUnauthorized, "you dont have the cookie!")
			}
			log.Panicln(err)
			return err
		}
		if cookie.Value == "kind of hash and...something" {
			return next(c)
		}
		return c.String(http.StatusUnauthorized, "you dont have right cookie")
	}
}
