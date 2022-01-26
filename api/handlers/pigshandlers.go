package handlers

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
)

type Pig struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func AddPigs(c echo.Context) error {
	pig := Pig{}
	defer c.Request().Body.Close()
	err := c.Bind(&pig)
	if err != nil {
		log.Printf("Failed Processing addPig request: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "")
	}
	log.Printf("this is your pig:%#v", pig)
	return c.String(http.StatusOK, "we got your pig!")

}
