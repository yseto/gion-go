package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"

	"github.com/yseto/gion-go/app/handler"
	"github.com/yseto/gion-go/config"
)

func main() {
	e := echo.New()

	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := sqlx.Open(cfg.DBDriverName, cfg.DBDataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()
	if dbConn.DriverName() == "sqlite3" {
		dbConn.Exec("PRAGMA foreign_keys = ON")
	}






	go func() {
		if err := e.Start(net.JoinHostPort(cfg.HTTPHost, cfg.HTTPPort)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
