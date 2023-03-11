package db_test

import (
	"errors"
	"log"
	"testing"

	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter11/db"
	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter11/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetMagazine(t *testing.T) {
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	t.Run("initial mag", func(t *testing.T) {
		ms := db.NewMagazineService(testDB, nil)
		em := ms.Upsert(db.Magazine{
			Name:   "New mag",
			Status: db.Available.String(),
		})
		log.Println(em)
		assert.NotNil(t, em)

		tests := map[string]struct {
			id      string
			want    db.Magazine
			wantErr error
		}{
			"existing mag": {id: em.ID, want: em},
			"no mag found": {id: "not-found", wantErr: db.ErrRecordNotFound},
			"empty id":     {id: "", wantErr: db.ErrRecordNotFound},
		}
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				m, err := ms.Get(tc.id)
				if tc.wantErr != nil {
					assert.Equal(t, tc.wantErr, err)
					assert.Nil(t, m)
					return
				}
				assert.Nil(t, err)
				assert.Equal(t, tc.want, *m)
			})
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		bs := db.NewMagazineService(testDB, nil)
		b, err := bs.Get("invalid-id")
		assert.Equal(t, db.ErrRecordNotFound, err)
		assert.Nil(t, b)
	})
}

func TestUpsertMag(t *testing.T) {
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	newMag := db.Magazine{
		Name:    "New mag",
		OwnerID: uuid.New().String(),
	}
	t.Run("new mag", func(t *testing.T) {
		ms := db.NewMagazineService(testDB, nil)
		m := ms.Upsert(newMag)
		assert.Equal(t, newMag.Name, m.Name)
		assert.Equal(t, newMag.OwnerID, m.OwnerID)
		assert.NotEmpty(t, m.ID)
		assert.Equal(t, db.Available.String(), m.Status)
	})

	t.Run("duplicate mag", func(t *testing.T) {
		ms := db.NewMagazineService(testDB, nil)
		m1 := ms.Upsert(newMag)
		m2 := ms.Upsert(m1)
		assert.Equal(t, m1, m2)
	})
}

func TestListMags(t *testing.T) {
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	t.Run("existing mags", func(t *testing.T) {
		ms := db.NewMagazineService(testDB, nil)
		em := ms.Upsert(db.Magazine{
			Name:   "Existing mag",
			Status: db.Available.String(),
		})
		mags, err := ms.List()
		require.Nil(t, err)
		assert.NotEmpty(t, mags)
		assert.Contains(t, mags, em)
	})

	t.Run("new mag", func(t *testing.T) {
		ms := db.NewMagazineService(testDB, nil)
		em := ms.Upsert(db.Magazine{
			Name:   "Existing mag",
			Status: db.Available.String(),
		})
		newMag := db.Magazine{
			Name:    "New mag",
			OwnerID: uuid.New().String(),
		}
		m := ms.Upsert(newMag)
		mags, err := ms.List()
		require.Nil(t, err)
		assert.NotEmpty(t, mags)
		assert.Contains(t, mags, em)
		assert.Contains(t, mags, m)
	})
}

func TestListMagsByUser(t *testing.T) {
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	t.Run("existing mag", func(t *testing.T) {
		ms := db.NewMagazineService(testDB, nil)
		em := ms.Upsert(db.Magazine{
			Name:    "Existing mag",
			Status:  db.Available.String(),
			OwnerID: uuid.New().String(),
		})
		mags, err := ms.ListByUser(em.OwnerID)
		require.Nil(t, err)
		assert.Equal(t, 1, len(mags))
		assert.Contains(t, mags, em)
	})

	t.Run("multiple mags", func(t *testing.T) {
		testDB, cleaner := db.OpenDB(t)
		defer cleaner()
		ms := db.NewMagazineService(testDB, nil)
		em := ms.Upsert(db.Magazine{
			Name:    "Existing mag",
			Status:  db.Available.String(),
			OwnerID: uuid.New().String(),
		})
		m := ms.Upsert(db.Magazine{
			Name:    "New mag",
			OwnerID: em.OwnerID,
		})
		mags, err := ms.ListByUser(m.OwnerID)
		require.Nil(t, err)
		assert.Equal(t, 2, len(mags))
		assert.Contains(t, mags, m)
		assert.Contains(t, mags, em)
	})

	t.Run("no mags for user", func(t *testing.T) {
		ms := db.NewMagazineService(testDB, nil)
		mags, err := ms.ListByUser(uuid.New().String())
		require.Nil(t, err)
		assert.Empty(t, mags)
	})
}

func TestSwapMagazine(t *testing.T) {
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	em := db.Magazine{
		Name:    "Existing mag",
		Status:  db.Available.String(),
		OwnerID: uuid.New().String(),
	}
	t.Run("existing mag", func(t *testing.T) {
		ps := mocks.NewPostingService(t)
		ms := db.NewMagazineService(testDB, ps)
		ps.On("NewMagazineOrder", mock.MatchedBy(func(m db.Magazine) bool {
			return m.ID == em.ID
		})).Return(nil).Once()
		em = ms.Upsert(em)
		newOwner := uuid.New().String()
		mag, err := ms.SwapMagazine(em.ID, newOwner)
		assert.NotNil(t, mag)
		assert.Nil(t, err)
		assert.Equal(t, em.ID, mag.ID)
		assert.Equal(t, newOwner, mag.OwnerID)
		assert.Equal(t, db.Swapped.String(), mag.Status)
	})

	t.Run("unknown mag", func(t *testing.T) {
		ps := mocks.NewPostingService(t)
		ms := db.NewMagazineService(testDB, ps)
		em = ms.Upsert(em)
		mag, err := ms.SwapMagazine(uuid.New().String(), uuid.New().String())
		assert.Nil(t, mag)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "no magazine found")
		ps.AssertNotCalled(t, "NewMagazineOrder", mock.AnythingOfType("db.Magazine"))
	})

	t.Run("empty list", func(t *testing.T) {
		ps := mocks.NewPostingService(t)
		ms := db.NewMagazineService(testDB, ps)
		mag, err := ms.SwapMagazine(uuid.New().String(), uuid.New().String())
		assert.Nil(t, mag)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "no magazine found")
		ps.AssertNotCalled(t, "NewMagazineOrder", mock.AnythingOfType("db.Magazine"))
	})

	t.Run("unavailable mag", func(t *testing.T) {
		ps := mocks.NewPostingService(t)
		ms := db.NewMagazineService(testDB, ps)
		em = ms.Upsert(em)
		ps.On("NewMagazineOrder", mock.MatchedBy(func(m db.Magazine) bool {
			return m.ID == em.ID
		})).Return(nil).Once()
		newOwner := uuid.New().String()
		mag, err := ms.SwapMagazine(em.ID, newOwner)
		assert.NotNil(t, mag)
		assert.Nil(t, err)
		assert.Equal(t, em.ID, mag.ID)
		assert.Equal(t, newOwner, mag.OwnerID)
		assert.Equal(t, db.Swapped.String(), mag.Status)
		mag, err = ms.SwapMagazine(em.ID, uuid.New().String())
		assert.Nil(t, mag)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "not available")
		ps.AssertExpectations(t)
	})

	t.Run("error posting", func(t *testing.T) {
		postingErr := errors.New("posting error")
		ps := mocks.NewPostingService(t)
		ms := db.NewMagazineService(testDB, ps)
		em = ms.Upsert(em)
		ps.On("NewMagazineOrder", mock.MatchedBy(func(m db.Magazine) bool {
			return m.ID == em.ID
		})).Return(postingErr).Once()
		newOwner := uuid.New().String()
		book, err := ms.SwapMagazine(em.ID, newOwner)
		assert.Nil(t, book)
		assert.Equal(t, postingErr, err)
		ps.AssertExpectations(t)
	})
}
