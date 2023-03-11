package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "embed"

	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter04/db"
	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter04/handlers"
)

//go:embed books.json
var booksFile []byte

//go:embed users.json
var usersFile []byte

func main() {
	port, ok := os.LookupEnv("BOOKSWAP_PORT")
	if !ok {
		log.Fatal("$BOOKSWAP_PORT not found")
	}
	books, users := importInitial()
	ps := db.NewPostingService()
	b := db.NewBookService(books, ps)
	u := db.NewUserService(users, b)
	h := handlers.NewHandler(b, u)

	router := handlers.ConfigureServer(h)
	log.Printf("Listening on :%s...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprint(":", port), router))
}

func importInitial() ([]db.Book, []db.User) {
	var books []db.Book
	var users []db.User

	err := json.Unmarshal(booksFile, &books)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(usersFile, &users)
	if err != nil {
		log.Fatal(err)
	}

	return books, users
}
