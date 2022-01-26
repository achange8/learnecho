package api

import (
	"github.com/achange8/learnecho/api/handlers"
	"github.com/labstack/echo"
)

func AdminGroup(g *echo.Group) {
	g.GET("/main", handlers.MainAdmin)
}
