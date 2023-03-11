package db_test

import (
	"testing"

	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter11/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewOrder(t *testing.T) {
	t.Run("book order", func(t *testing.T) {
		ps := db.NewPostingService()
		b := db.Book{
			ID: uuid.New().String(),
		}
		err := ps.NewBookOrder(b)
		assert.Nil(t, err)
	})
	t.Run("mag order", func(t *testing.T) {
		ps := db.NewPostingService()
		m := db.Magazine{
			ID: uuid.New().String(),
		}
		err := ps.NewMagazineOrder(m)
		assert.Nil(t, err)
	})
}
