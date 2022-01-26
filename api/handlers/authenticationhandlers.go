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
		//create jwt access token
		Access_Token, err := createAccessToken()
		if err != nil {
			log.Println("Err Creating JWT Access_Token!", err)
			return c.String(http.StatusInternalServerError, "some thing wrong")
		}
		JWTCookie := new(http.Cookie)

		JWTCookie.Name = "JWT_Access_CCookie"
		JWTCookie.Value = Access_Token
		JWTCookie.Expires = time.Now().Add(30 * time.Minute)
		JWTCookie.HttpOnly = true

		c.SetCookie(JWTCookie)

		//create jwt refresh token
		Refresh_Token, err := createRefreshToken()
		if err != nil {
			log.Println("Err Creating JWT Refresh_Token!", err)
			return c.String(http.StatusInternalServerError, "some thing wrong")
		}
		JWTRefreshCookie := new(http.Cookie)

		JWTRefreshCookie.Name = "JWT_Access_CCookie"
		JWTRefreshCookie.Value = Refresh_Token
		JWTRefreshCookie.Expires = time.Now().Add(15 * time.Minute)
		JWTRefreshCookie.HttpOnly = true

		c.SetCookie(JWTRefreshCookie)

		return c.JSON(http.StatusOK, map[string]string{
			"message":      "You were logged in!",
			"Access_Token": Access_Token,
		})
	}
	return c.JSON(http.StatusUnauthorized, "Worng imformation!")
}

////////////////////create Tokens///////////////////////////
func createAccessToken() (string, error) {
	claims := JwtClaims{
		"osh",
		jwt.StandardClaims{
			Id:        "main_user_id",
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		},
	}
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	token, err := rawToken.SignedString([]byte("mySecret"))
	if err != nil {
		return "", err
	}
	return token, nil
}

func createRefreshToken() (string, error) {
	claims := JwtClaims{
		"osh",
		jwt.StandardClaims{
			Id:        "main_user_id",
			ExpiresAt: time.Now().Add(24 * 14 * time.Hour).Unix(),
		},
	}
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	token, err := rawToken.SignedString([]byte("mySecret"))
	if err != nil {
		return "", err
	}
	return token, nil
}
