package db_test

import (
	"errors"
	"testing"

	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter07/db"
	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter07/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetBook(t *testing.T) {
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	t.Run("initial books", func(t *testing.T) {
		bs := db.NewBookRepository(testDB, nil)
		eb := bs.Upsert(db.Book{
			Name:   "New Book",
			Status: db.Available.String(),
		})
		assert.NotNil(t, eb)

		tests := map[string]struct {
			book    db.Book
			want    db.Book
			wantErr error
		}{
			"existing book": {book: db.Book{ID: eb.ID}, want: eb},
			"no book found": {book: db.Book{ID: "not-found"}, wantErr: db.ErrRecordNotFound},
			"empty id":      {book: db.Book{}, wantErr: db.ErrRecordNotFound},
		}
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				err := bs.Get(&tc.book)
				if tc.wantErr != nil {
					assert.Equal(t, tc.wantErr, err)
					return
				}
				assert.Nil(t, err)
				assert.Equal(t, tc.want, eb)
			})
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		bs := db.NewBookRepository(testDB, nil)
		ib := db.Book{
			ID: "invalid-id",
		}
		err := bs.Get(&ib)
		assert.Equal(t, db.ErrRecordNotFound, err)
	})
}

func TestUpsertBook(t *testing.T) {
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	newBook := db.Book{
		Name:    "New book",
		OwnerID: uuid.New().String(),
	}
	t.Run("new book", func(t *testing.T) {
		bs := db.NewBookRepository(testDB, nil)
		b := bs.Upsert(newBook)
		assert.Equal(t, newBook.Name, b.Name)
		assert.Equal(t, newBook.OwnerID, b.OwnerID)
		assert.NotEmpty(t, b.ID)
		assert.Equal(t, db.Available.String(), b.Status)
	})

	t.Run("duplicate book", func(t *testing.T) {
		bs := db.NewBookRepository(testDB, nil)
		b1 := bs.Upsert(newBook)
		b2 := bs.Upsert(b1)
		assert.Equal(t, b1, b2)
	})
}

func TestListBooks(t *testing.T) {
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	t.Run("existing books", func(t *testing.T) {
		bs := db.NewBookRepository(testDB, nil)
		eb := bs.Upsert(db.Book{
			Name:   "Existing book",
			Status: db.Available.String(),
		})
		books, err := bs.List()
		require.Nil(t, err)
		assert.NotEmpty(t, books)
		assert.Contains(t, books, eb)
	})

	t.Run("new book", func(t *testing.T) {
		bs := db.NewBookRepository(testDB, nil)
		eb := bs.Upsert(db.Book{
			Name:   "Existing book",
			Status: db.Available.String(),
		})
		newBook := db.Book{
			Name:    "New book",
			OwnerID: uuid.New().String(),
		}
		b := bs.Upsert(newBook)
		books, err := bs.List()
		require.Nil(t, err)
		assert.NotEmpty(t, books)
		assert.Contains(t, books, eb)
		assert.Contains(t, books, b)
	})
}

func TestListBooksByUser(t *testing.T) {
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	t.Run("existing book", func(t *testing.T) {
		bs := db.NewBookRepository(testDB, nil)
		eb := bs.Upsert(db.Book{
			Name:    "Existing book",
			Status:  db.Available.String(),
			OwnerID: uuid.New().String(),
		})
		books, err := bs.ListByUser(eb.OwnerID)
		require.Nil(t, err)
		assert.Equal(t, 1, len(books))
		assert.Contains(t, books, eb)
	})

	t.Run("multiple books", func(t *testing.T) {
		testDB, cleaner := db.OpenDB(t)
		defer cleaner()
		bs := db.NewBookRepository(testDB, nil)
		eb := bs.Upsert(db.Book{
			Name:    "Existing book",
			Status:  db.Available.String(),
			OwnerID: uuid.New().String(),
		})
		b := bs.Upsert(db.Book{
			Name:    "New book",
			OwnerID: eb.OwnerID,
		})
		books, err := bs.ListByUser(b.OwnerID)
		require.Nil(t, err)
		assert.Equal(t, 2, len(books))
		assert.Contains(t, books, b)
		assert.Contains(t, books, eb)
	})

	t.Run("no books for user", func(t *testing.T) {
		bs := db.NewBookRepository(testDB, nil)
		books, err := bs.ListByUser(uuid.New().String())
		require.Nil(t, err)
		assert.Empty(t, books)
	})
}

func TestSwapBook(t *testing.T) {
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	eb := db.Book{
		Name:    "Existing book",
		Status:  db.Available.String(),
		OwnerID: uuid.New().String(),
	}
	t.Run("existing book", func(t *testing.T) {
		ps := mocks.NewPostingService(t)
		bs := db.NewBookRepository(testDB, ps)
		ps.On("NewOrder", mock.MatchedBy(func(b db.Book) bool {
			return b.ID == eb.ID
		})).Return(nil).Once()
		eb = bs.Upsert(eb)
		newOwner := uuid.New().String()
		book, err := bs.SwapBook(eb.ID, newOwner)
		assert.NotNil(t, book)
		assert.Nil(t, err)
		assert.Equal(t, eb.ID, book.ID)
		assert.Equal(t, newOwner, book.OwnerID)
		assert.Equal(t, db.Swapped.String(), book.Status)
		ps.AssertExpectations(t)
	})

	t.Run("unknown book", func(t *testing.T) {
		ps := mocks.NewPostingService(t)
		bs := db.NewBookRepository(testDB, ps)
		eb = bs.Upsert(eb)
		book, err := bs.SwapBook(uuid.New().String(), uuid.New().String())
		assert.Nil(t, book)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "no book found")
		ps.AssertNotCalled(t, "NewOrder", mock.AnythingOfType("db.Book"))
	})

	t.Run("empty list", func(t *testing.T) {
		ps := mocks.NewPostingService(t)
		bs := db.NewBookRepository(testDB, ps)
		book, err := bs.SwapBook(uuid.New().String(), uuid.New().String())
		assert.Nil(t, book)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "no book found")
		ps.AssertNotCalled(t, "NewOrder", mock.AnythingOfType("db.Book"))
	})

	t.Run("unavailable book", func(t *testing.T) {
		ps := mocks.NewPostingService(t)
		bs := db.NewBookRepository(testDB, ps)
		eb = bs.Upsert(eb)
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
		bs := db.NewBookRepository(testDB, ps)
		eb = bs.Upsert(eb)
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
