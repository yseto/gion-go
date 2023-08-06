package main

//go:generate go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.13.0
//go:generate oapi-codegen -package handler -generate echo-server,strict-server -o handler/server.gen.go openapi.yaml
//go:generate oapi-codegen -package handler -generate spec -o handler/spec.gen.go openapi.yaml
//go:generate oapi-codegen -package handler -generate types -o handler/types.gen.go openapi.yaml

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	oapiMiddleware "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
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

	swagger, _ := handler.GetSwagger()
	swagger.Servers = nil

	e.Use(
		oapiMiddleware.OapiRequestValidatorWithOptions(swagger,
			&oapiMiddleware.Options{
				Options: openapi3filter.Options{
					AuthenticationFunc: handler.NewAuthenticator(),
				},
			}),

		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				ctx := handler.NewDBContext(c.Request().Context(), dbConn)

				// into a value http.Request.Context from echo.Context
				if userid := c.Get(handler.SessionContextKey); userid != nil {
					ctx = handler.NewUserContext(ctx, userid.(uint64))
				}
				c.SetRequest(c.Request().WithContext(ctx))
				return next(c)
			}
		})
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	handler.RegisterHandlers(e, handler.NewStrictHandler(handler.NewApiServer(), nil))

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
