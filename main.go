package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Product struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
type Dog struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
type Pig struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "hello from the web side!")
}

func getproducts(c echo.Context) error {
	P_Name := c.QueryParam("name")
	P_Type := c.QueryParam("type")

	dataType := c.Param("data")
	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("Product name is %s\n type is %s \n", P_Name, P_Type))
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

func addproduct(c echo.Context) error {
	product := Product{}
	defer c.Request().Body.Close()
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed reading the request body : %s", err)
		return c.String(http.StatusInternalServerError, "")
	}
	err = json.Unmarshal(b, &product)
	if err != nil {
		log.Printf("Failed Unmarshaling in addProducts: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}
	log.Printf("this is your product:%#v", product)
	return c.String(http.StatusOK, "we got your product!")
}

func addDog(c echo.Context) error {
	dog := Dog{}
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&dog)
	if err != nil {
		log.Printf("Failed Processing addDog request: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "")
	}
	log.Printf("this is your dog:%#v", dog)
	return c.String(http.StatusOK, "we got your dog!")
}

func addPigs(c echo.Context) error {
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

func mainAdmin(g echo.Context) error {
	return g.String(http.StatusOK, "horay you are on secret admin main page!")
}

func mainCookie(g echo.Context) error {
	return g.String(http.StatusOK, "you are on secret cookie main page!")
}

func login(c echo.Context) error {
	userID := c.QueryParam("userID")
	password := c.QueryParam("password")

	//check userid and password against DB after hashing the password
	if userID == "osh" && password == "1234" {
		cookie := new(http.Cookie)

		cookie.Name = "sessionID"
		cookie.Value = "kind of hash and...something"
		cookie.Expires = time.Now().Add(24 * time.Hour)

		c.SetCookie(cookie)
		return c.String(http.StatusOK, "You were logged in!")
	}
	return c.String(http.StatusUnauthorized, "Worng imformation!")
}

///////////////////////////middleware///////////////////////
func serverHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "osh/1.0")
		c.Response().Header().Set("echo.HeaderServer", "osh/1.2")

		return next(c)
	}
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

/////////////////////////////////////////////////////////////////
func main() {
	fmt.Println("Welcome osh server with echo!")

	e := echo.New()
	e.Use(serverHeader)
	adminGroup := e.Group("/admin")

	cookieGroup := e.Group("/cookie")

	adminGroup.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${host} ${path} ${latency_human}` + "\n",
	}))

	adminGroup.Use(middleware.BasicAuth(func(userID, Password string, c echo.Context) (bool, error) {
		if userID == "osh" && Password == "1234" {
			return true, nil
		}
		return false, nil
	}))

	cookieGroup.Use(checkCookie)
	adminGroup.GET("/main", mainAdmin)
	e.GET("/", hello)
	e.GET("/login", login)
	cookieGroup.GET("/main", mainCookie)

	e.GET("/products/:id", getproducts)
	e.POST("/products", addproduct)
	e.POST("/dogs", addDog)
	e.POST("/pigs", addPigs)
	e.Logger.Fatal(e.Start(":8080"))
}
