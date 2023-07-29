package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/yseto/gion-go/config"
)

func Open(cfg *config.Config) (*sqlx.DB, error) {

	dbConn, err := sqlx.Open(cfg.DBDriverName, cfg.DBDataSourceName)
	if err != nil {
		return nil, err
	}
	if dbConn.DriverName() == "sqlite3" {
		dbConn.Exec("PRAGMA foreign_keys = ON")
	}
	return dbConn, nil
}
