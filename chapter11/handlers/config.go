package handlers

import (
	"net/http"
	"os"

	_ "net/http/pprof"

	"github.com/gorilla/mux"
)

// ConfigureServer configures the routes of this server and binds handler functions to them
func ConfigureServer(handler *Handler) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Methods("GET").Path("/").Handler(http.HandlerFunc(handler.Index))
	router.Methods("GET").Path("/books").Handler(http.HandlerFunc(handler.ListBooks))
	router.Methods("POST").Path("/users").Handler(http.HandlerFunc(handler.UserUpsert))
	router.Methods("GET").Path("/users/{id}/books").Handler(http.HandlerFunc(handler.ListUserByID_Books))
	router.Methods("POST").Path("/books/{id}").Handler(http.HandlerFunc(handler.SwapBook))
	router.Methods("POST").Path("/books").Handler(http.HandlerFunc(handler.BookUpsert))
	router.Methods("GET").Path("/users/{id}/magazines").Handler(http.HandlerFunc(handler.ListUserByID_Magazines))
	router.Methods("GET").Path("/magazines").Handler(http.HandlerFunc(handler.ListMagazines))
	router.Methods("POST").Path("/magazines").Handler(http.HandlerFunc(handler.MagazineUpsert))
	router.Methods("POST").Path("/magazines/{id}").Handler(http.HandlerFunc(handler.SwapMagazine))

	if os.Getenv("DEBUG") != "" {
		router.PathPrefix("/debug/pprof/").
			Handler(http.DefaultServeMux)
	}
	return router
}
