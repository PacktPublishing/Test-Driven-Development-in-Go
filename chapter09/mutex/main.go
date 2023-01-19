package main

import (
	"fmt"
	"sync"
)

// Thread safe LIFO Stack implementation
type Stack struct {
	lock sync.Mutex
	data []string
}

func NewStack() *Stack {
	return &Stack{
		data: make([]string, 0),
	}
}

// Push adds the given element to the end of the list
func (s *Stack) Push(el string) {
	defer s.lock.Unlock()
	s.lock.Lock()
	s.data = append(s.data, el)
}

// Pop removes and returns the last element from the list,
// or an error if the list is empty.
func (s *Stack) Pop() (*string, error) {
	defer s.lock.Unlock()
	s.lock.Lock()
	if len(s.data) == 0 {
		return nil, fmt.Errorf("stack is empty")
	}
	last := s.data[len(s.data)-1]
	s.data = s.data[0 : len(s.data)-1]
	return &last, nil
}
