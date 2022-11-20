package db_test

import (
	"errors"
	"testing"

	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter05/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
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

// TODO: Test the other methods of the the BookService
