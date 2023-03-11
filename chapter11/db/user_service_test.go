package db_test

import (
	"testing"

	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter11/db"
	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter11/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUser(t *testing.T) {
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	t.Run("existing user", func(t *testing.T) {
		eb := db.Book{
			Name:    "Existing book",
			OwnerID: uuid.New().String(),
		}
		em := db.Magazine{
			Name:    "Existing mag",
			OwnerID: uuid.New().String(),
		}
		bs := mocks.NewBookOperationsService(t)
		ms := mocks.NewMagazineOperationsService(t)
		us := db.NewUserService(testDB, bs, ms)
		eu, err := us.Upsert(db.User{
			Name: "Existing user",
		})
		require.Nil(t, err)
		bs.On("ListByUser", eu.ID).Return([]db.Book{eb}, nil).Once()
		ms.On("ListByUser", eu.ID).Return([]db.Magazine{em}, nil).Once()
		userProfile, err := us.Get(eu.ID)
		assert.Nil(t, err)
		assert.Equal(t, eu, userProfile.User)
		assert.Equal(t, 1, len(userProfile.Books))
		assert.Contains(t, userProfile.Books, eb)
		bs.AssertExpectations(t)
		ms.AssertExpectations(t)
	})
	t.Run("invalid users", func(t *testing.T) {
		us := db.NewUserService(testDB, nil, nil)
		tests := map[string]struct {
			id      string
			wantErr string
		}{
			"no user found": {id: "not-found", wantErr: "no user found"},
			"empty id":      {id: "", wantErr: "no user found"},
		}
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				userProfile, err := us.Get(tc.id)
				assert.Contains(t, err.Error(), tc.wantErr)
				assert.Nil(t, userProfile)
			})
		}
	})
}

func TestUpsertUser(t *testing.T) {
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	bs := mocks.NewBookOperationsService(t)
	ms := mocks.NewMagazineOperationsService(t)
	us := db.NewUserService(testDB, bs, ms)
	newUser := db.User{
		Name: "New user",
	}
	user, err := us.Upsert(newUser)
	require.Nil(t, err)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, newUser.Name, user.Name)
	bs.AssertNotCalled(t, "ListByUser")
	ms.AssertNotCalled(t, "ListByUser")
}

func TestExistsUser(t *testing.T) {
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	bs := mocks.NewBookOperationsService(t)
	ms := mocks.NewMagazineOperationsService(t)
	t.Run("existing user", func(t *testing.T) {
		us := db.NewUserService(testDB, bs, ms)
		eu, err := us.Upsert(db.User{
			Name: "Existing user",
		})
		require.Nil(t, err)
		err = us.Exists(eu.ID)
		require.Nil(t, err)
	})
	t.Run("invalid ID user", func(t *testing.T) {
		us := db.NewUserService(testDB, bs, ms)
		err := us.Exists(uuid.New().String())
		require.NotNil(t, err)
		assert.Contains(t, err.Error(), "no user found")
	})
}
