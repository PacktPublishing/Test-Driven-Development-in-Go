package db

import (
	"log"
)

// PostingService interface wraps around external posting functionality.
type PostingService interface {
	NewBookOrder(b Book) error
	NewMagazineOrder(m Magazine) error
}

// StubbedPostingService is a concrete mock of the external PostingService.
type StubbedPostingService struct{}

// NewPostingService initialises the PostingService.
func NewPostingService() PostingService {
	return &StubbedPostingService{}
}

// NewBookOrder creates a book order and sends it to the posting servivce for posting.
func (sps *StubbedPostingService) NewBookOrder(b Book) error {
	log.Printf("STUBBED POSTING SERVICE: book %s posted: %v", b.ID, b)
	return nil
}

// NewMagazineOrder creates a book order and sends it to the posting servivce for posting.
func (sps *StubbedPostingService) NewMagazineOrder(m Magazine) error {
	log.Printf("STUBBED POSTING SERVICE: magazine %s posted: %v", m.ID, m)
	return nil
}
