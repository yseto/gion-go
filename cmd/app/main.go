package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/lestrrat-go/accesslog"

	"github.com/yseto/gion-go/app/handler"
	"github.com/yseto/gion-go/config"
	"github.com/yseto/gion-go/db"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := db.Open(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	al := accesslog.New().Logger(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	svr := handler.New()
	h := handler.Handler(handler.NewStrictHandler(svr, nil))

	mw, err := handler.CreateMiddleware(cfg.JwtSignedKeyBin)
	if err != nil {
		log.Fatalln("error creating middleware:", err)
	}
	h = mw(h)

	mw2 := handler.CreateMiddlewareEmptyContext(dbConn)
	h = mw2(h)

	s := &http.Server{
		Handler: al.Wrap(h),
		Addr:    net.JoinHostPort(cfg.HTTPHost, cfg.HTTPPort),
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
		close(idleConnsClosed)
	}()

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("shutting down the server")
	}

	<-idleConnsClosed
}
