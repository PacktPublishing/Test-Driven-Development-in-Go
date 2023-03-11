package errors_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorsVerification(t *testing.T) {
	t.Run("simple custom error", func(t *testing.T) {
		wantMsg := "Something went wrong!"
		err := errors.New(wantMsg)
		gotMsg := err.Error()
		assert.Equal(t, wantMsg, gotMsg)
	})
	t.Run("formatted custom error", func(t *testing.T) {
		input := 4
		wantMsg := fmt.Sprintf("Input %d cannot be even.", input)
		err := checkOdd(input)
		gotMsg := err.Error()
		assert.Equal(t, wantMsg, gotMsg)
	})
	t.Run("contains custom error", func(t *testing.T) {
		input := 4
		err := checkOdd(input)
		gotMsg := err.Error()
		assert.Contains(t, gotMsg, fmt.Sprint(input))
		assert.Contains(t, gotMsg, "even")
	})
	t.Run("custom error type", func(t *testing.T) {
		input := 4
		wantErr := &evenNumberError{
			input: input,
		}
		err := checkOdd(input)
		var gotErr *evenNumberError
		require.True(t, errors.As(err, &gotErr))
		assert.Equal(t, wantErr, gotErr)
	})
}

type evenNumberError struct {
	input int
}

func (e *evenNumberError) Error() string {
	return fmt.Sprintf("Input %d cannot be even.", e.input)
}

func checkOdd(input int) error {
	if input%2 == 0 {
		return &evenNumberError{
			input: input,
		}
	}
	return nil
}
