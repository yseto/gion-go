package handler

import (
	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/yseto/gion-go/db"
)

type CustomContext struct {
	echo.Context
	Conn *sqlx.DB
}

func (c *CustomContext) DBCommon() db.DB {
	return db.New(c.Conn)
}

func (c *CustomContext) DBUser() db.UserScopedDB {
	user := c.Get(SessionContextKey).(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	return db.NewUserScopedDB(c.Conn, claims.UserID)
}

const SessionContextKey = "user"
const JwtSignedKey = "secret"

func JWTConfig() middleware.JWTConfig {
	return middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte(JwtSignedKey),
		ContextKey: SessionContextKey,
	}
}
