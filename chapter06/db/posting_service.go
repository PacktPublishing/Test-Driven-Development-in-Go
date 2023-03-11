package db

import (
	"log"
)

// PostingService interface wraps around external posting functionality.
type PostingService interface {
	NewOrder(b Book) error
}

// StubbedPostingService is a concrete mock of the external PostingService.
type StubbedPostingService struct{}

// NewPostingService initialises the PostingService.
func NewPostingService() PostingService {
	return &StubbedPostingService{}
}

// NewOrder creates a new order and sends it to the posting servivce for posting.
func (sps *StubbedPostingService) NewOrder(b Book) error {
	log.Printf("STUBBED POSTING SERVICE: book %s posted: %v", b.ID, b)
	return nil
}
