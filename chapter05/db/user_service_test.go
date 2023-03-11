package db_test

import (
	"testing"

	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter05/db"
	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter05/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUser(t *testing.T) {
	t.Run("existing user", func(t *testing.T) {
		eb := db.Book{
			Name:    "New book",
			OwnerID: uuid.New().String(),
		}
		eu := db.User{
			ID:   uuid.New().String(),
			Name: "Existing user",
		}
		bs := mocks.NewBookOperationsService(t)
		us := db.NewUserService([]db.User{eu}, bs)
		bs.On("ListByUser", eu.ID).Return([]db.Book{eb}).Once()
		user, books, err := us.Get(eu.ID)
		assert.Nil(t, err)
		assert.Equal(t, eu, *user)
		assert.Equal(t, 1, len(books))
		assert.Contains(t, books, eb)
		bs.AssertExpectations(t)
	})
	t.Run("invalid users", func(t *testing.T) {
		eu := db.User{
			ID:   uuid.New().String(),
			Name: "Existing user",
		}
		us := db.NewUserService([]db.User{eu}, nil)
		tests := map[string]struct {
			id      string
			wantErr string
		}{
			"no user found": {id: "not-found", wantErr: "no user found"},
			"empty id":      {id: "", wantErr: "no user found"},
		}
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				user, books, err := us.Get(tc.id)
				assert.Contains(t, err.Error(), tc.wantErr)
				assert.Nil(t, user)
				assert.Empty(t, books)
			})
		}
	})

	t.Run("empty users", func(t *testing.T) {
		us := db.NewUserService([]db.User{}, nil)
		user, books, err := us.Get(uuid.New().String())
		assert.Nil(t, user)
		assert.Empty(t, books)
		assert.Contains(t, err.Error(), "no user found")
	})
}

func TestUpsertUser(t *testing.T) {
	bs := mocks.NewBookOperationsService(t)
	us := db.NewUserService([]db.User{}, bs)
	newUser := db.User{
		Name: "New user",
	}
	user, err := us.Upsert(newUser)
	require.Nil(t, err)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, newUser.Name, user.Name)
	bs.AssertNotCalled(t, "ListByUser")
}

func TestExistsUser(t *testing.T) {
	eu := db.User{
		ID:   uuid.New().String(),
		Name: "Existing user",
	}
	bs := mocks.NewBookOperationsService(t)
	t.Run("existing user", func(t *testing.T) {
		us := db.NewUserService([]db.User{eu}, bs)
		err := us.Exists(eu.ID)
		require.Nil(t, err)
	})
	t.Run("invalid ID user", func(t *testing.T) {
		us := db.NewUserService([]db.User{eu}, bs)
		err := us.Exists(uuid.New().String())
		require.NotNil(t, err)
		assert.Contains(t, err.Error(), "no user found")
	})
	t.Run("empty users", func(t *testing.T) {
		us := db.NewUserService([]db.User{}, bs)
		err := us.Exists(uuid.New().String())
		require.NotNil(t, err)
		assert.Contains(t, err.Error(), "no user found")
	})
}
