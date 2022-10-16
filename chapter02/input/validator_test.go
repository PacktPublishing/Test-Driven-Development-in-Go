package input_test

import (
	"testing"

	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter02/input"
)

func TestCheckInput(t *testing.T) {
	validOperators := []string{"+"}
	t.Run("valid case", func(t *testing.T) {
		v := setup(t, validOperators)
		err := v.CheckInput(validOperators[0], []float64{2.5, 3.5})
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("invalid operator", func(t *testing.T) {
		v := setup(t, validOperators)
		err := v.CheckInput("-", []float64{2.5, 3.5})
		if err == nil {
			t.Fatal(err)
		}
	})
	t.Run("invalid operand count", func(t *testing.T) {
		v := setup(t, validOperators)
		err := v.CheckInput("-", []float64{2.5})
		if err == nil {
			t.Fatal(err)
		}
	})
}

func setup(t *testing.T, validOps []string) *input.Validator {
	t.Helper()
	return input.NewValidator(2, validOps)
}
