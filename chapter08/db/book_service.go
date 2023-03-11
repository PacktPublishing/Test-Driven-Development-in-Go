package db

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Book contains all the fields for representing a book.
type Book struct {
	ID      string `json:"id" gorm:"primaryKey"`
	Name    string `json:"name"`
	Author  string `json:"author"`
	OwnerID string `json:"owner_id"`
	Status  string `json:"status"`
}

// BookService contains all the functionality and dependencies for managing books.
type BookService struct {
	DB *gorm.DB
	ps PostingService
}

// NewBookService initialises a BookService given its dependencies.
func NewBookService(db *gorm.DB, ps PostingService) *BookService {
	return &BookService{
		DB: db,
		ps: ps,
	}
}

// Get returns a given book or error if none exists.
func (bs *BookService) Get(id string) (*Book, error) {
	var b Book
	if r := bs.DB.Where("id = ?", id).First(&b); r.Error != nil {
		return nil, r.Error
	}

	return &b, nil
}

// Upsert creates or updates a book.
func (bs *BookService) Upsert(b Book) Book {
	var eb Book
	if r := bs.DB.Where("id = ?", b.ID).First(&eb); r.Error != nil {
		b.ID = uuid.NewString()
		b.Status = Available.String()
	}
	bs.DB.Save(&b)
	return b
}

// List returns the list of available books.
func (bs *BookService) List() ([]Book, error) {
	var items []Book
	if result := bs.DB.Where("status = ?", Available.String()).Find(&items); result.Error != nil {
		return nil, result.Error
	}

	return items, nil
}

// ListByUser returns the list of books for a given user.
func (bs *BookService) ListByUser(userID string) ([]Book, error) {
	var items []Book
	if result := bs.DB.Where("owner_id = ?", userID).Find(&items); result.Error != nil {
		return nil, result.Error
	}

	return items, nil
}

// SwapBook checks whether a book is available and, if possible, marks it as swapped.
func (bs *BookService) SwapBook(bookID, userID string) (*Book, error) {
	var b Book
	if r := bs.DB.Where("id = ?", bookID).First(&b); r.Error != nil {
		return nil, fmt.Errorf("no book found for id %s:%v", bookID, r.Error)
	}
	if b.Status != Available.String() {
		return nil, fmt.Errorf("book %s is not available for swapping", bookID)
	}
	b.OwnerID = userID
	b.Status = Swapped.String()
	sb := bs.Upsert(b)
	if err := bs.ps.NewOrder(sb); err != nil {
		return nil, err
	}

	return &sb, nil
}
