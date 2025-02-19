package handler

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/yseto/gion-go/db"

	myContext "github.com/yseto/gion-go/app/handler/context"
)

type myContextKey struct{}

const (
	myUserKey   = "userid"
	myDBConnKey = "dbConn"
)

func emptyUserContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, myContextKey{}, myContext.New())
}

func newUserContext(ctx context.Context, userid uint64) {
	ctx.Value(myContextKey{}).(*myContext.Context).Set(myUserKey, userid)
}

func userFromContext(ctx context.Context) uint64 {
	return ctx.Value(myContextKey{}).(*myContext.Context).Get(myUserKey).(uint64)
}

func NewDBContext(ctx context.Context, h *sqlx.DB) {
	ctx.Value(myContextKey{}).(*myContext.Context).Set(myDBConnKey, h)
}

func dbFromContext(ctx context.Context) *sqlx.DB {
	return ctx.Value(myContextKey{}).(*myContext.Context).Get(myDBConnKey).(*sqlx.DB)
}

func DBCommonFromContext(ctx context.Context) db.DB {
	return db.New(dbFromContext(ctx))
}

func DBUserFromContext(ctx context.Context) db.UserScopedDB {
	return db.NewUserScopedDB(dbFromContext(ctx), userFromContext(ctx))
}
