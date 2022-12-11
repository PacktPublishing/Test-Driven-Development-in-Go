package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter07/db"
	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter07/handlers"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var postgresURL = "postgres://adelinasimion:postgres@localhost:5432/bookswap?sslmode=disable"

func main() {
	m, err := migrate.New("file://db/migrations", postgresURL)
	if err != nil {
		log.Fatal(err)
	}	
	if err := m.Up(); err != migrate.ErrNoChange {
		log.Fatal(err)
	}
	// defer func() {
	// 	m.Down()
	// }()
	dbConn, err := gorm.Open(postgres.Open(postgresURL), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	ps := db.NewPostingService()
	b := db.NewBookService(dbConn, ps)
	u := db.NewUserService(dbConn, b)
	h := handlers.NewHandler(b, u)

	router := handlers.ConfigureServer(h)
	fmt.Println("Listening on localhost:3000...")
	log.Fatal(http.ListenAndServe("localhost:3000", router))
}
