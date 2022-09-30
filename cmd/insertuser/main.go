package main

import (
	"flag"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"

	"github.com/yseto/gion-go/config"
	"github.com/yseto/gion-go/db"
)

func main() {
	username := flag.String("u", "username", "set username")
	password := flag.String("p", "password", "set password")
	flag.Parse()

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

	digest, err := bcrypt.GenerateFromPassword([]byte(*password), 8)
	if err != nil {
		log.Fatal(err)
	}

	dbc := db.New(dbConn)
	if err := dbc.InsertUser(*username, string(digest)); err != nil {
		log.Fatal(err)
	}
	log.Printf("create user : %s\n", *username)

	user, err := dbc.UserByName(*username)
	if err != nil {
		log.Fatal(err)
	}

	tx := db.NewUserScopedDB(dbConn, user.ID).MustBegin()
	if err := tx.InsertCategory("no title"); err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}
