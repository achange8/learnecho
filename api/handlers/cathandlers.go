package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

type Cat struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func GetCats(c echo.Context) error {
	P_Name := c.QueryParam("name")
	P_Type := c.QueryParam("type")

	dataType := c.Param("data")
	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("Cat name is %s\n type is %s \n", P_Name, P_Type))
	}
	if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"name": P_Name,
			"type": P_Type,
		})
	}
	return c.JSON(http.StatusBadRequest, map[string]string{
		"error": "You need to let us know if you want json or string data",
	})
}

func AddCats(c echo.Context) error {
	Cat := Cat{}
	defer c.Request().Body.Close()
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed reading the request body : %s", err)
		return c.String(http.StatusInternalServerError, "")
	}
	err = json.Unmarshal(b, &Cat)
	if err != nil {
		log.Printf("Failed Unmarshaling in addCats: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}
	log.Printf("this is your Cat:%#v", Cat)
	return c.String(http.StatusOK, "we got your Cat!")
}
