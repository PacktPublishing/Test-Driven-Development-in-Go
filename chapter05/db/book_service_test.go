package db_test

import (
	"errors"
	"testing"

	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter05/db"
	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter05/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetBook(t *testing.T) {
	t.Run("initial books", func(t *testing.T) {
		eb := db.Book{
			ID:     uuid.New().String(),
			Name:   "Existing book",
			Status: db.Available.String(),
		}
		bs := db.NewBookService([]db.Book{eb}, nil)

		tests := map[string]struct {
			id      string
			want    db.Book
			wantErr error
		}{
			"existing book": {id: eb.ID, want: eb},
			"no book found": {id: "not-found", wantErr: errors.New("no book found")},
			"empty id":      {id: "", wantErr: errors.New("no book found")},
		}
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				b, err := bs.Get(tc.id)
				if tc.wantErr != nil {
					assert.Equal(t, tc.wantErr, err)
					assert.Nil(t, b)
					return
				}
				assert.Nil(t, err)
				assert.Equal(t, tc.want, *b)
			})
		}
	})

	t.Run("empty books", func(t *testing.T) {
		bs := db.NewBookService([]db.Book{}, nil)
		b, err := bs.Get("id")
		assert.Equal(t, errors.New("no book found"), err)
		assert.Nil(t, b)
	})
}

func TestUpsertBook(t *testing.T) {
	newBook := db.Book{
		Name:    "New book",
		OwnerID: uuid.New().String(),
	}
	t.Run("new book", func(t *testing.T) {
		bs := db.NewBookService([]db.Book{}, nil)
		b := bs.Upsert(newBook)
		assert.Equal(t, newBook.Name, b.Name)
		assert.Equal(t, newBook.OwnerID, b.OwnerID)
		assert.NotEmpty(t, b.ID)
		assert.Equal(t, db.Available.String(), b.Status)
	})

	t.Run("duplicate book", func(t *testing.T) {
		bs := db.NewBookService([]db.Book{}, nil)
		b1 := bs.Upsert(newBook)
		b2 := bs.Upsert(b1)
		assert.Equal(t, b1, b2)
	})
}

func TestListBooks(t *testing.T) {
	eb := db.Book{
		ID:     uuid.New().String(),
		Name:   "Existing book",
		Status: db.Available.String(),
	}
	t.Run("existing books", func(t *testing.T) {
		bs := db.NewBookService([]db.Book{eb}, nil)
		books := bs.List()
		assert.Equal(t, 1, len(books))
		assert.Contains(t, books, eb)
	})

	t.Run("new book", func(t *testing.T) {
		bs := db.NewBookService([]db.Book{eb}, nil)
		newBook := db.Book{
			Name:    "New book",
			OwnerID: uuid.New().String(),
		}
		b := bs.Upsert(newBook)
		books := bs.List()
		assert.Equal(t, 2, len(books))
		assert.Contains(t, books, eb)
		assert.Contains(t, books, b)
	})
	t.Run("empty books", func(t *testing.T) {
		bs := db.NewBookService([]db.Book{}, nil)
		books := bs.List()
		assert.Empty(t, books)
	})
}

func TestListBooksByUser(t *testing.T) {
	eb := db.Book{
		ID:      uuid.New().String(),
		Name:    "Existing book",
		Status:  db.Available.String(),
		OwnerID: uuid.New().String(),
	}
	t.Run("existing book", func(t *testing.T) {
		bs := db.NewBookService([]db.Book{eb}, nil)
		books := bs.ListByUser(eb.OwnerID)
		assert.Equal(t, 1, len(books))
		assert.Contains(t, books, eb)
	})

	t.Run("new book", func(t *testing.T) {
		bs := db.NewBookService([]db.Book{eb}, nil)
		newBook := db.Book{
			Name:    "New book",
			OwnerID: uuid.New().String(),
		}
		b := bs.Upsert(newBook)
		books := bs.ListByUser(b.OwnerID)
		assert.Equal(t, 1, len(books))
		assert.Contains(t, books, b)
	})

	t.Run("multiple books", func(t *testing.T) {
		bs := db.NewBookService([]db.Book{eb}, nil)
		newBook := db.Book{
			Name:    "New book",
			OwnerID: eb.OwnerID,
		}
		b := bs.Upsert(newBook)
		books := bs.ListByUser(b.OwnerID)
		assert.Equal(t, 2, len(books))
		assert.Contains(t, books, b)
		assert.Contains(t, books, eb)
	})

	t.Run("no books for user", func(t *testing.T) {
		bs := db.NewBookService([]db.Book{eb}, nil)
		books := bs.ListByUser(uuid.New().String())
		assert.Empty(t, books)
	})
}

func TestSwapBook(t *testing.T) {
	eb := db.Book{
		ID:      uuid.New().String(),
		Name:    "Existing book",
		Status:  db.Available.String(),
		OwnerID: uuid.New().String(),
	}
	t.Run("existing book", func(t *testing.T) {
		ps := mocks.NewPostingService(t)
		bs := db.NewBookService([]db.Book{eb}, ps)
		ps.On("NewOrder", mock.MatchedBy(func(b db.Book) bool {
			return b.ID == eb.ID
		})).Return(nil).Once()
		newOwner := uuid.New().String()
		book, err := bs.SwapBook(eb.ID, newOwner)
		assert.NotNil(t, book)
		assert.Nil(t, err)
		assert.Equal(t, eb.ID, book.ID)
		assert.Equal(t, newOwner, book.OwnerID)
		assert.Equal(t, db.Swapped.String(), book.Status)
		ps.AssertExpectations(t)
	})

	t.Run("unknown book id", func(t *testing.T) {
		ps := mocks.NewPostingService(t)
		bs := db.NewBookService([]db.Book{eb}, ps)
		book, err := bs.SwapBook(uuid.New().String(), uuid.New().String())
		assert.Nil(t, book)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "no book found")
		ps.AssertNotCalled(t, "NewOrder", mock.AnythingOfType("db.Book"))
	})

	t.Run("empty list", func(t *testing.T) {
		ps := mocks.NewPostingService(t)
		bs := db.NewBookService([]db.Book{}, ps)
		book, err := bs.SwapBook(uuid.New().String(), uuid.New().String())
		assert.Nil(t, book)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "no book found")
		ps.AssertNotCalled(t, "NewOrder", mock.AnythingOfType("db.Book"))
	})

	t.Run("unavailable book", func(t *testing.T) {
		ps := mocks.NewPostingService(t)
		bs := db.NewBookService([]db.Book{eb}, ps)
		ps.On("NewOrder", mock.MatchedBy(func(b db.Book) bool {
			return b.ID == eb.ID
		})).Return(nil).Once()
		newOwner := uuid.New().String()
		book, err := bs.SwapBook(eb.ID, newOwner)
		assert.NotNil(t, book)
		assert.Nil(t, err)
		assert.Equal(t, eb.ID, book.ID)
		assert.Equal(t, newOwner, book.OwnerID)
		assert.Equal(t, db.Swapped.String(), book.Status)
		book, err = bs.SwapBook(eb.ID, uuid.New().String())
		assert.Nil(t, book)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "not available")
		ps.AssertExpectations(t)
	})

	t.Run("error posting", func(t *testing.T) {
		postingErr := errors.New("posting error")
		ps := mocks.NewPostingService(t)
		bs := db.NewBookService([]db.Book{eb}, ps)
		ps.On("NewOrder", mock.MatchedBy(func(b db.Book) bool {
			return b.ID == eb.ID
		})).Return(postingErr).Once()
		newOwner := uuid.New().String()
		book, err := bs.SwapBook(eb.ID, newOwner)
		assert.Nil(t, book)
		assert.Equal(t, postingErr, err)
		ps.AssertExpectations(t)
	})
}
