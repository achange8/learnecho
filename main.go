package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
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

type JwtClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
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

func mainJWT(c echo.Context) error {
	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	log.Println("User name:", claims["name"], "User ID:", claims["jti"])

	return c.JSON(http.StatusOK, " you are on the secrt page in jwt")
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
		//TODO: create jwt token
		token, err := createJwtToken(userID)
		if err != nil {
			log.Println("Err Creating JWT token!", err)
			return c.String(http.StatusInternalServerError, "some thing wrong")
		}
		JWTCookie := new(http.Cookie)

		JWTCookie.Name = "JWTCookie"
		JWTCookie.Value = token
		JWTCookie.Expires = time.Now().Add(24 * time.Hour)

		c.SetCookie(JWTCookie)

		return c.JSON(http.StatusOK, map[string]string{
			"message": "You were logged in!",
			"token":   token,
		})
	}
	return c.JSON(http.StatusUnauthorized, "Worng imformation!")
}

func createJwtToken(userID string) (string, error) {
	claims := JwtClaims{
		userID,
		jwt.StandardClaims{
			Id:        "main_user_id",
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	token, err := rawToken.SignedString([]byte("mySecret"))
	if err != nil {
		return "", err
	}
	return token, nil
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

	jwtGroup := e.Group("/jwt")

	adminGroup.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${host} ${path} ${latency_human}` + "\n",
	}))

	adminGroup.Use(middleware.BasicAuth(func(userID, Password string, c echo.Context) (bool, error) {
		if userID == "osh" && Password == "1234" {
			return true, nil
		}
		return false, nil
	}))

	jwtGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte("mySecret"),
		TokenLookup:   "cookie:JWTCookie",
		// - "query:<name>"
		// - "cookie:<name>"
	}))

	cookieGroup.Use(checkCookie)

	adminGroup.GET("/main", mainAdmin)

	cookieGroup.GET("/main", mainCookie)

	jwtGroup.GET("/main", mainJWT)

	e.GET("/", hello)
	e.GET("/login", login)
	e.GET("/products/:id", getproducts)
	e.POST("/products", addproduct)
	e.POST("/dogs", addDog)
	e.POST("/pigs", addPigs)
	e.Logger.Fatal(e.Start(":8080"))
}
