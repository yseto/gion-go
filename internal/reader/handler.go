package reader

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/yseto/gion-go/internal/config"
)

func CreateHandler(cfg *config.Config, dbConn *sqlx.DB) http.Handler {
	svr := New(cfg.JwtSignedKeyBin)
	h := Handler(NewStrictHandler(svr, nil))

	mw, err := CreateMiddleware(cfg.JwtSignedKeyBin)
	if err != nil {
		log.Fatalln("error creating middleware:", err)
	}
	h = mw(h)

	mw2 := CreateMiddlewareEmptyContext(dbConn)
	h = mw2(h)

	return h
}
