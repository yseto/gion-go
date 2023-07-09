package handler

import (
	"github.com/golang-jwt/jwt"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/yseto/gion-go/db"
)

type userID struct{}

func NewUserContext(ctx context.Context, userid uint64) context.Context {
	return context.WithValue(ctx, userID{}, userid)
}

func userFromContext(ctx context.Context) uint64 {
	return ctx.Value(userID{}).(uint64)
}

type dbConnKey struct{}

func NewDBContext(ctx context.Context, h *sqlx.DB) context.Context {
	return context.WithValue(ctx, dbConnKey{}, h)
}

func dbFromContext(ctx context.Context) *sqlx.DB {
	return ctx.Value(dbConnKey{}).(*sqlx.DB)
}

func DBCommonFromContext(ctx context.Context) db.DB {
	return db.New(dbFromContext(ctx))
}

func DBUserFromContext(ctx context.Context) db.UserScopedDB {
	return db.NewUserScopedDB(dbFromContext(ctx), userFromContext(ctx))
}

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
