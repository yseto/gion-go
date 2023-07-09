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
