package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	oapiMiddleware "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/yseto/gion-go/app/handler"
	"github.com/yseto/gion-go/config"
	"github.com/yseto/gion-go/db"
)

func main() {
	e := echo.New()

	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := db.Open(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	swagger, _ := handler.GetSwagger()
	swagger.Servers = nil

	e.Use(
		oapiMiddleware.OapiRequestValidatorWithOptions(swagger,
			&oapiMiddleware.Options{
				Options: openapi3filter.Options{
					AuthenticationFunc: handler.NewAuthenticator(cfg.JwtSignedKeyBin),
				},
				ErrorHandler: func(c echo.Context, err *echo.HTTPError) error {
					if message, ok := err.Message.(string); ok {
						if strings.HasPrefix(message, "security requirements failed") {
							c.Response().WriteHeader(401)
							return nil
						}
					}
					return err
				},
			}),

		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				ctx := handler.NewDBContext(c.Request().Context(), dbConn)

				ctx = config.NewContext(ctx, cfg)

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
