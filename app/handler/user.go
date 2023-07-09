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
