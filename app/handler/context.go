package handler

import (
	"context"

	"github.com/jmoiron/sqlx"

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
