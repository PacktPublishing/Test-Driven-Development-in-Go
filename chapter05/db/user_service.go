package db

import (
	"fmt"

	"github.com/google/uuid"
)

// User contains all the user fields.
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	PostCode string `json:"post_code"`
	Country  string `json:"country"`
}

// UserService has all the dependencies required for managing users.
type UserService struct {
	users map[string]User
	bs    *BookService
}

// NewUserService initialises the UserService.
func NewUserService(initial []User, bs *BookService) *UserService {
	users := make(map[string]User)
	for _, u := range initial {
		users[u.ID] = u
	}
	return &UserService{
		users: users,
		bs:    bs,
	}
}

// Get returns a given user or error if none exists.
func (us *UserService) Get(id string) (*User, []Book, error) {
	u, ok := us.users[id]
	if !ok {
		return nil, nil, fmt.Errorf("no user found for id %s", id)
	}
	books := us.bs.ListByUser(id)

	return &u, books, nil
}

// Exists returns whether a given user exists and returns an error if none found.
func (us *UserService) Exists(id string) error {
	_, ok := us.users[id]
	if !ok {
		return fmt.Errorf("no user found for id %s", id)
	}

	return nil
}

// Upsert creates or updates a new order.
func (us *UserService) Upsert(u User) (User, error) {
	u.ID = uuid.NewString()
	us.users[u.ID] = u

	return u, nil
}
