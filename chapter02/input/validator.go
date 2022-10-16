package input

import (
	"fmt"
)

type Validator struct {
	expectedLength int
	validOperators []string
}

// NewValidator creates a ready to use Validator instance
func NewValidator(expLen int, validOps []string) *Validator {
	return &Validator{
		expectedLength: expLen,
		validOperators: validOps,
	}
}

// CheckInput validates the operator and operands
func (v *Validator) CheckInput(operator string, operands []float64) error {
	opLen := len(operands)
	if opLen != v.expectedLength {
		return fmt.Errorf("unexpected operands length: got %d, want %d", opLen, v.expectedLength)
	}
	return v.checkOperator(operator)
}

// checkOperator validates the operator is supported
func (v *Validator) checkOperator(operator string) error {
	for _, o := range v.validOperators {
		if o == operator {
			return nil
		}
	}

	return fmt.Errorf("invalid operator:%s", operator)
}
