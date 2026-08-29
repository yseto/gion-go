package engine

import "fmt"

const (
	PostgreSQL = "pgx"
	MySQL      = "mysql"
	SQLite     = "sqlite3"
)

func NormalizeDriverName(driverName string) string {
	switch driverName {
	case "postgres":
		return PostgreSQL
	case "sqlite3":
		return SQLite
	case "mysql":
		return MySQL
	}
	panic(fmt.Sprintf("%s is invalid", driverName))
}
