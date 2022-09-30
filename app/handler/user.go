package handler

import (
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type jwtCustomClaims struct {
	UserID uint64 `json:"user_id"`
	jwt.StandardClaims
}

// https://echo.labstack.com/middleware/jwt/
// https://echo.labstack.com/cookbook/jwt/
func Login(c echo.Context) error {
	username := c.FormValue("id")
	password := c.FormValue("password")

	db := c.(*CustomContext).DBCommon()

	user, err := db.UserByName(username)
	if err != nil {
		c.Response().Header().Set("WWW-Authenticate", `Bearer realm="need token" error="invalid_token"`)
		return echo.ErrUnauthorized
	}

	if check := bcrypt.CompareHashAndPassword([]byte(user.Digest), []byte(password)); check != nil {
		c.Response().Header().Set("WWW-Authenticate", `Bearer realm="need token" error="invalid_token"`)
		return echo.ErrUnauthorized
	}

	if err := db.UpdateUserLastLogin(user.ID); err != nil {
		fmt.Println(err)
	}

	claims := &jwtCustomClaims{
		user.ID,
		jwt.StandardClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	signedToken, err := token.SignedString([]byte(JwtSignedKey))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"autoseen": user.UserProfile.AutoSeen,
		"token":    signedToken,
	})
}

func Logout(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{})
}

type updatePasswordResult struct {
	Result string `json:"result"`
}

func UpdatePassword(c echo.Context) error {
	passwordOld := c.FormValue("password_old")
	password := c.FormValue("password")
	passwordCheck := c.FormValue("passwordc")

	if password != passwordCheck || utf8.RuneCountInString(password) < 8 {
		return c.JSON(http.StatusOK, updatePasswordResult{"error"})
	}

	db := c.(*CustomContext).DBUser()
	user, err := db.User()
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if check := bcrypt.CompareHashAndPassword([]byte(user.Digest), []byte(passwordOld)); check != nil {
		return c.JSON(http.StatusOK, updatePasswordResult{"unmatch now password"})
	}

	newDigest, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if db.UpdateUserDigest(string(newDigest)) != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusOK, updatePasswordResult{"update password"})
}
