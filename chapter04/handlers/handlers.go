package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter04/db"
	"github.com/gorilla/mux"
)

// Handler contains the handler and all its dependencies.
type Handler struct {
	bs *db.BookService
	us *db.UserService
}

// NewHandler initialises a new handler, given dependencies.
func NewHandler(bs *db.BookService, us *db.UserService) *Handler {
	return &Handler{
		bs: bs,
		us: us,
	}
}

// Index is invoked by HTTP GET /.
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	// Send an HTTP status & a hardcoded message
	resp := &Response{
		Message: "Welcome to the BookSwap service!",
		Books:   h.bs.List(),
	}
	writeResponse(w, http.StatusOK, resp)
}

// ListBooks is invoked by HTTP GET /books.
func (h *Handler) ListBooks(w http.ResponseWriter, r *http.Request) {
	// Send an HTTP status & the list of books
	writeResponse(w, http.StatusOK, &Response{
		Books: h.bs.List(),
	})
}

// UserUpsert is invoked by HTTP POST /users.
func (h *Handler) UserUpsert(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := readRequestBody(r)
	// Handle any errors & write an error HTTP status & response
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, &Response{
			Error: fmt.Errorf("invalid user body:%v", err).Error(),
		})
		return
	}

	// Initialize a user to unmarshal request body into
	var user db.User
	if err := json.Unmarshal(body, &user); err != nil {
		writeResponse(w, http.StatusUnprocessableEntity, &Response{
			Error: fmt.Errorf("invalid user body:%v", err).Error(),
		})
		return
	}

	// Call the repository method corresponding to the operation
	user, err = h.us.Upsert(user)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, &Response{
			Error: err.Error(),
		})
		return
	}

	// Send an HTTP success status & the return value from the repo
	writeResponse(w, http.StatusOK, &Response{
		User: &user,
	})
}

// ListUserByID is invoked by HTTP GET /users/{id}.
func (h *Handler) ListUserByID(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]
	user, books, err := h.us.Get(userID)
	if err != nil {
		writeResponse(w, http.StatusNotFound, &Response{
			Error: err.Error(),
		})
		return
	}

	// Send an HTTP success status & the return value from the repo
	writeResponse(w, http.StatusOK, &Response{
		User:  user,
		Books: books,
	})
}

// SwapBook is invoked by POST /users/{id}
func (h *Handler) SwapBook(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]
	if err := h.us.Exists(userID); err != nil {
		writeResponse(w, http.StatusBadRequest, &Response{
			Error: err.Error(),
		})
		return
	}

	bookID := r.URL.Query().Get("book")
	if bookID == "" {
		resp := &Response{
			Error: fmt.Errorf("no bookID provided for book swap").Error(),
		}
		writeResponse(w, http.StatusBadRequest, resp)

	}
	book, err := h.bs.SwapBook(bookID, userID)
	if err != nil {
		writeResponse(w, http.StatusNotFound, &Response{
			Error: err.Error(),
		})
		return
	}

	writeResponse(w, http.StatusNotFound, &Response{
		Books: []db.Book{*book},
	})
}

// BookUpsert is invoked by HTTP POST /books.
func (h *Handler) BookUpsert(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := readRequestBody(r)
	// Handle any errors & write an error HTTP status & response
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, &Response{
			Error: fmt.Errorf("invalid book body:%v", err).Error(),
		})
		return
	}

	// Initialize a book to unmarshal request body into
	var book db.Book
	if err := json.Unmarshal(body, &book); err != nil {
		writeResponse(w, http.StatusUnprocessableEntity, &Response{
			Error: fmt.Errorf("invalid book body:%v", err).Error(),
		})
		return
	}
	if err := h.us.Exists(book.OwnerID); err != nil {
		writeResponse(w, http.StatusBadRequest, &Response{
			Error: err.Error(),
		})
		return
	}

	// Call the repository method corresponding to the operation
	book = h.bs.Upsert(book)
	// Send an HTTP success status & the return value from the repo
	writeResponse(w, http.StatusOK, &Response{
		Books: []db.Book{book},
	})
}

// readRequestBody is a helper method that
// allows to read a request body and return any errors.
func readRequestBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return []byte{}, err
	}
	if err := r.Body.Close(); err != nil {
		return []byte{}, err
	}
	return body, err
}
