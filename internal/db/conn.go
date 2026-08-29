package db

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/yseto/gion-go/internal/config"
	"github.com/yseto/gion-go/internal/db/engine"
)

func Open(cfg *config.Config) (*sqlx.DB, error) {
	dbConn, err := sqlx.Open(engine.NormalizeDriverName(cfg.DBDriverName), cfg.DBDataSourceName)
	if err != nil {
		return nil, err
	}
	if dbConn.DriverName() == engine.SQLite {
		dbConn.Exec("PRAGMA foreign_keys = ON")
	}
	return dbConn, nil
}
