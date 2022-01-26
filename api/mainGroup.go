package api

import (
	"github.com/achange8/learnecho/api/handlers"
	"github.com/labstack/echo"
)

func MainGroup(e *echo.Echo) {
	e.GET("/login", handlers.Login)
	e.GET("/hello", handlers.Hello)
	e.GET("/cats/:id", handlers.GetCats)

	e.POST("/cats", handlers.AddCats)
	e.POST("/dogs", handlers.AddDog)
	e.POST("/pigs", handlers.AddPigs)

}
