package db

// BooksStatus contains the different types of Book status.
type BookStatus int

const (
	Available BookStatus = iota
	Swapped
)

func (o BookStatus) String() string {
	return [...]string{"AVAILABLE", "SWAPPED"}[o]
}
