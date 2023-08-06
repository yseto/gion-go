package main

import (
	"flag"
	"log"

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

	dbConn, err := db.Open(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

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
