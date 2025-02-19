package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	// TODO

	idleConnsClosed := make(chan struct{})

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}
		close(idleConnsClosed)
	}()

	if err := e.Start(net.JoinHostPort(cfg.HTTPHost, cfg.HTTPPort)); err != nil && err != http.ErrServerClosed {
		e.Logger.Fatal("shutting down the server")
	}

	<-idleConnsClosed
}
