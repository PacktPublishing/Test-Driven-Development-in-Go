package db

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Magazine contains all the fields for representing a magazine.
type Magazine struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	IssueNumber int    `json:"issue_number"`
	OwnerID     string `json:"owner_id"`
	Status      string `json:"status"`
}

// MagazineService contains all the functionality and dependencies for managing magazines.
type MagazineService struct {
	DB *gorm.DB
	ps PostingService
}

// NewMagazineService initialises a MagazineService given its dependencies.
func NewMagazineService(db *gorm.DB, ps PostingService) *MagazineService {
	return &MagazineService{
		DB: db,
		ps: ps,
	}
}

// Get returns a given magazine or error if none exists.
func (ms *MagazineService) Get(id string) (*Magazine, error) {
	var m Magazine
	if r := ms.DB.Where("id = ?", id).First(&m); r.Error != nil {
		return nil, r.Error
	}

	return &m, nil
}

// Upsert creates or updates a magazine.
func (ms *MagazineService) Upsert(m Magazine) Magazine {
	var em Magazine
	if r := ms.DB.Where("id = ?", m.ID).First(&em); r.Error != nil {
		m.ID = uuid.NewString()
		m.Status = Available.String()
	}
	ms.DB.Save(&m)
	return m
}

// List returns the list of available magazines.
func (ms *MagazineService) List() ([]Magazine, error) {
	var items []Magazine
	if result := ms.DB.Where("status = ?", Available.String()).Find(&items); result.Error != nil {
		return nil, result.Error
	}

	return items, nil
}

// ListByUser returns the list of magazines for a given user.
func (ms *MagazineService) ListByUser(userID string) ([]Magazine, error) {
	var items []Magazine
	if result := ms.DB.Where("owner_id = ?", userID).Find(&items); result.Error != nil {
		return nil, result.Error
	}

	return items, nil
}

// SwapMagazine checks whether a magazine is available and, if possible, marks it as swapped.
func (ms *MagazineService) SwapMagazine(magID, userID string) (*Magazine, error) {
	var m Magazine
	if r := ms.DB.Where("id = ?", magID).First(&m); r.Error != nil {
		return nil, fmt.Errorf("no magazine found for id %s:%v", magID, r.Error)
	}
	if m.Status != Available.String() {
		return nil, fmt.Errorf("mag %s is not available for swapping", magID)
	}
	m.OwnerID = userID
	m.Status = Swapped.String()
	sm := ms.Upsert(m)
	if err := ms.ps.NewMagazineOrder(sm); err != nil {
		return nil, err
	}

	return &sm, nil
}
