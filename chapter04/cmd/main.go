package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "embed"

	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter04/db"
	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter04/handlers"
)

//go:embed books.json
var booksFile []byte

//go:embed users.json
var usersFile []byte

func main() {
	books, users := importInitial()
	ps := db.NewPostingService()
	b := db.NewBookService(books, ps)
	u := db.NewUserService(users, b)
	h := handlers.NewHandler(b, u)

	router := handlers.ConfigureServer(h)
	fmt.Println("Listening on localhost:3000...")
	log.Fatal(http.ListenAndServe("localhost:3000", router))
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
