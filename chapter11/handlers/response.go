package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter11/db"
)
type ResponseItemType interface {
	db.Book | db.Magazine
}

// Response contains all the response types of our handlers.
type Response[T ResponseItemType] struct {
	Message string   `json:"message,omitempty"`
	Error   string   `json:"error,omitempty"`
	Items   []T      `json:"items,omitempty"`
	User    *db.User `json:"user,omitempty"`
}

// writeResponse is a helper method that allows to write the HTTP status & response
func writeResponse[T ResponseItemType](w http.ResponseWriter, status int, resp *Response[T]) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if status != http.StatusOK {
		w.WriteHeader(status)
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Fprintf(w, "error encoding resp %v:%s", resp, err)
	}
}
