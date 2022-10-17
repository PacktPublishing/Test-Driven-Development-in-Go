package table

import (
	"errors"
	"fmt"
)

func Divide(x, y int8) (*string, error) {
	if y == 0 {
		return nil, errors.New("cannot divide by 0")
	}

	r := float64(x) / float64(y)
	result := fmt.Sprintf("%.2f", r)
	return &result, nil
}
