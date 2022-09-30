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

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &handler.CustomContext{Context: c, Conn: dbConn}
			return next(cc)
		}
	})
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/api/login", handler.Login, handler.CheckXHR())
	e.POST("/api/logout", handler.Logout, handler.CheckXHR())

	e.Static("/", "public")

	r := e.Group("/api/")
	r.Use(middleware.JWTWithConfig(handler.JWTConfig()))
	r.Use(handler.CheckXHR())

	r.POST("register_category", handler.RegisterCategory)
	r.POST("register_subscription", handler.RegisterSubscription)
	r.POST("examine_subscription", handler.ExamineSubscription)
	r.POST("delete_subscription", handler.DeleteSubscription)
	r.POST("change_subscription", handler.ChangeSubscription)
	r.POST("profile", handler.Profile)
	r.POST("set_profile", handler.UpdateProfile)
	r.POST("categories", handler.Categories)
	r.POST("category_with_count", handler.CategoryAndUnreadEntryCount)
	r.POST("unread_entry", handler.UnreadEntry)
	r.POST("set_asread", handler.SetAsread)
	r.POST("subscriptions", handler.Subscriptions)
	r.POST("pinned_items", handler.PinnedItems)
	r.POST("set_pin", handler.SetPin)
	r.POST("remove_all_pin", handler.RemoveAllPin)
	r.POST("update_password", handler.UpdatePassword)
	r.POST("opml_export", handler.OpmlExport)
	r.POST("opml_import", handler.OpmlImport)

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
