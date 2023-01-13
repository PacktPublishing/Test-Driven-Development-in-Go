package handlers

import (
	"net/http"

	_ "net/http/pprof"

	"github.com/gorilla/mux"
)

// ConfigureServer configures the routes of this server and binds handler functions to them
func ConfigureServer(handler *Handler) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Methods("GET").Path("/").Handler(http.HandlerFunc(handler.Index))
	router.Methods("GET").Path("/books").Handler(http.HandlerFunc(handler.ListBooks))
	router.Methods("POST").Path("/users").Handler(http.HandlerFunc(handler.UserUpsert))
	router.Methods("GET").Path("/users/{id}").Handler(http.HandlerFunc(handler.ListUserByID))
	router.Methods("POST").Path("/books/{id}").Handler(http.HandlerFunc(handler.SwapBook))
	router.Methods("POST").Path("/books").Handler(http.HandlerFunc(handler.BookUpsert))
	router.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)

	return router
}
