package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter05/db"
)

// Response contains all the response types of our handlers.
type Response struct {
	Message string    `json:"message,omitempty"`
	Error   string    `json:"error,omitempty"`
	Books   []db.Book `json:"books,omitempty"`
	User    *db.User  `json:"user,omitempty"`
}

// writeResponse is a helper method that allows to write the HTTP status & response
func writeResponse(w http.ResponseWriter, status int, resp *Response) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if status != http.StatusOK {
		w.WriteHeader(status)
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Fprintf(w, "error encoding resp %v:%s", resp, err)
	}
}
