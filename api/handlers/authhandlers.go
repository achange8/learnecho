package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
)

type JwtClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func Login(c echo.Context) error {
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
		token, err := createJwtToken()
		if err != nil {
			log.Println("Err Creating JWT token!", err)
			return c.String(http.StatusInternalServerError, "some thing wrong")
		}
		JWTCookie := new(http.Cookie)

		JWTCookie.Name = "JWTCookie"
		JWTCookie.Value = token
		JWTCookie.Expires = time.Now().Add(24 * time.Hour)
		JWTCookie.HttpOnly = true

		c.SetCookie(JWTCookie)

		return c.JSON(http.StatusOK, map[string]string{
			"message": "You were logged in!",
			"token":   token,
		})
	}
	return c.JSON(http.StatusUnauthorized, "Worng imformation!")
}

func createJwtToken() (string, error) {
	claims := JwtClaims{
		"osh",
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
