package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter11/db"
	"github.com/gorilla/mux"
)

// Handler contains the handler and all its dependencies.
type Handler struct {
	bs *db.BookService
	us *db.UserService
	ms *db.MagazineService
}

// NewHandler initialises a new handler, given dependencies.
func NewHandler(bs *db.BookService, us *db.UserService, ms *db.MagazineService) *Handler {
	return &Handler{
		bs: bs,
		us: us,
		ms: ms,
	}
}

// Index is invoked by HTTP GET /.
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	books, err := h.bs.List()
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, &Response[db.Book]{
			Error: err.Error(),
		})
		return
	}

	// Send an HTTP status & a hardcoded message
	resp := &Response[db.Book]{
		Message: "Welcome to the BookSwap service!",
		Items:   books,
	}
	writeResponse(w, http.StatusOK, resp)
}

// ListBooks is invoked by HTTP GET /books.
func (h *Handler) ListBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.bs.List()
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, &Response[db.Book]{
			Error: err.Error(),
		})
		return
	}

	// Send an HTTP status & the list of books
	writeResponse(w, http.StatusOK, &Response[db.Book]{
		Items: books,
	})
}

// ListMagazines is invoked by HTTP GET /magazines.
func (h *Handler) ListMagazines(w http.ResponseWriter, r *http.Request) {
	mags, err := h.ms.List()
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, &Response[db.Magazine]{
			Error: err.Error(),
		})
		return
	}

	// Send an HTTP status & the list of mags
	writeResponse(w, http.StatusOK, &Response[db.Magazine]{
		Items: mags,
	})
}

// UserUpsert is invoked by HTTP POST /users.
func (h *Handler) UserUpsert(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := readRequestBody(r)
	// Handle any errors & write an error HTTP status & response
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, &Response[db.Book]{
			Error: fmt.Errorf("invalid user body:%v", err).Error(),
		})
		return
	}

	// Initialize a user to unmarshal request body into
	var user db.User
	if err := json.Unmarshal(body, &user); err != nil {
		writeResponse(w, http.StatusUnprocessableEntity, &Response[db.Book]{
			Error: fmt.Errorf("invalid user body:%v", err).Error(),
		})
		return
	}

	// Call the repository method corresponding to the operation
	user, err = h.us.Upsert(user)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, &Response[db.Book]{
			Error: err.Error(),
		})
		return
	}

	// Send an HTTP success status & the return value from the repo
	writeResponse(w, http.StatusOK, &Response[db.Book]{
		User: &user,
	})
}

// ListUserByID_Books is invoked by HTTP GET /users/{id}/books.
func (h *Handler) ListUserByID_Books(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]
	userProfile, err := h.us.Get(userID)
	if err != nil {
		writeResponse(w, http.StatusNotFound, &Response[db.Book]{
			Error: err.Error(),
		})
		return
	}

	// Send an HTTP success status & the return value from the repo
	writeResponse(w, http.StatusOK, &Response[db.Book]{
		User:  &userProfile.User,
		Items: userProfile.Books,
	})
}

// ListUserByID_Magazines is invoked by HTTP GET /users/{id}/magazines.
func (h *Handler) ListUserByID_Magazines(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]
	userProfile, err := h.us.Get(userID)
	if err != nil {
		writeResponse(w, http.StatusNotFound, &Response[db.Magazine]{
			Error: err.Error(),
		})
		return
	}

	// Send an HTTP success status & the return value from the repo
	writeResponse(w, http.StatusOK, &Response[db.Magazine]{
		User:  &userProfile.User,
		Items: userProfile.Magazines,
	})
}

// SwapBook is invoked by POST /books/{id}
func (h *Handler) SwapBook(w http.ResponseWriter, r *http.Request) {
	bookID := mux.Vars(r)["id"]
	userID := r.URL.Query().Get("user")
	if err := h.us.Exists(userID); err != nil {
		writeResponse(w, http.StatusBadRequest, &Response[db.Book]{
			Error: err.Error(),
		})
		return
	}
	_, err := h.bs.SwapBook(bookID, userID)
	if err != nil {
		writeResponse(w, http.StatusNotFound, &Response[db.Book]{
			Error: err.Error(),
		})
		return
	}

	userProfile, err := h.us.Get(userID)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, &Response[db.Book]{
			Error: err.Error(),
		})
		return
	}

	writeResponse(w, http.StatusOK, &Response[db.Book]{
		User:  &userProfile.User,
		Items: userProfile.Books,
	})
}

// SwapMagazine is invoked by POST /magazines/{id}
func (h *Handler) SwapMagazine(w http.ResponseWriter, r *http.Request) {
	magID := mux.Vars(r)["id"]
	userID := r.URL.Query().Get("user")
	if err := h.us.Exists(userID); err != nil {
		writeResponse(w, http.StatusBadRequest, &Response[db.Book]{
			Error: err.Error(),
		})
		return
	}
	_, err := h.ms.SwapMagazine(magID, userID)
	if err != nil {
		writeResponse(w, http.StatusNotFound, &Response[db.Magazine]{
			Error: err.Error(),
		})
		return
	}

	userProfile, err := h.us.Get(userID)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, &Response[db.Magazine]{
			Error: err.Error(),
		})
		return
	}

	writeResponse(w, http.StatusOK, &Response[db.Magazine]{
		User:  &userProfile.User,
		Items: userProfile.Magazines,
	})
}

// BookUpsert is invoked by HTTP POST /books.
func (h *Handler) BookUpsert(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := readRequestBody(r)
	// Handle any errors & write an error HTTP status & response
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, &Response[db.Book]{
			Error: fmt.Errorf("invalid book body:%v", err).Error(),
		})
		return
	}

	// Initialize a book to unmarshal request body into
	var book db.Book
	if err := json.Unmarshal(body, &book); err != nil {
		writeResponse(w, http.StatusUnprocessableEntity, &Response[db.Book]{
			Error: fmt.Errorf("invalid book body:%v", err).Error(),
		})
		return
	}
	if err := h.us.Exists(book.OwnerID); err != nil {
		writeResponse(w, http.StatusBadRequest, &Response[db.Book]{
			Error: err.Error(),
		})
		return
	}

	// Call the repository method corresponding to the operation
	updatedBook := h.bs.Upsert(book)
	// Send an HTTP success status & the return value from the repo
	writeResponse(w, http.StatusOK, &Response[db.Book]{
		Items: []db.Book{updatedBook},
	})
}

// MagazineUpsert is invoked by HTTP POST /magazines.
func (h *Handler) MagazineUpsert(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := readRequestBody(r)
	// Handle any errors & write an error HTTP status & response
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, &Response[db.Magazine]{
			Error: fmt.Errorf("invalid mag body:%v", err).Error(),
		})
		return
	}

	// Initialize a mag to unmarshal request body into
	var mag db.Magazine
	if err := json.Unmarshal(body, &mag); err != nil {
		writeResponse(w, http.StatusUnprocessableEntity, &Response[db.Magazine]{
			Error: fmt.Errorf("invalid mag body:%v", err).Error(),
		})
		return
	}
	if err := h.us.Exists(mag.OwnerID); err != nil {
		writeResponse(w, http.StatusBadRequest, &Response[db.Magazine]{
			Error: err.Error(),
		})
		return
	}

	// Call the repository method corresponding to the operation
	updatedMag := h.ms.Upsert(mag)
	// Send an HTTP success status & the return value from the repo
	writeResponse(w, http.StatusOK, &Response[db.Magazine]{
		Items: []db.Magazine{updatedMag},
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
