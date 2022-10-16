package calculator_test

import (
	"testing"

	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter03/calculator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetExpectedLength(t *testing.T) {
	// Arrange
	e := calculator.NewEngine()

	// Act
	opLen := e.GetNumOperands()

	//Assert
	assert.Equal(t, 2, opLen)
}

func TestGetValidOperations(t *testing.T) {
	// Arrange
	e := calculator.NewEngine()

	// Act
	ops := e.GetValidOperators()

	// Assert
	assert.Equal(t, 1, len(ops))
	assert.Contains(t, ops, "+")
}

func TestAdd(t *testing.T) {
	// Arrange
	e := calculator.NewEngine()

	t.Run("positive input", func(t *testing.T) {
		// Arrange
		x, y := 2.5, 3.5
		want := 6.0

		// Act
		result := e.Add(x, y)

		// Assert
		assert.Equal(t, want, result)
	})

	t.Run("negative input", func(t *testing.T) {
		// Arrange
		x, y := -2.5, -3.5
		want := -6.0

		// Act
		result := e.Add(x, y)

		// Assert
		assert.Equal(t, want, result)
	})
}

func TestProcessOperation(t *testing.T) {
	t.Run("valid operation", func(t *testing.T) {
		// Arrange
		e := calculator.NewEngine()
		op := calculator.Operation{
			Expression: "2 + 3",
			Operator:   "+",
			Operands:   []float64{2, 3},
		}
		expectedResult := "5.0"

		// Act
		result, err := e.ProcessOperation(op)

		// Assert
		require.Nil(t, err)
		require.NotNil(t, result)
		assert.Contains(t, *result, expectedResult)
		assert.Contains(t, *result, op.Expression)
	})

	t.Run("invalid operation", func(t *testing.T) {
		// Arrange
		e := calculator.NewEngine()
		op := calculator.Operation{
			Expression: "2 % 3",
			Operator:   "%",
			Operands:   []float64{2, 3},
		}

		// Act
		result, err := e.ProcessOperation(op)

		// Assert
		require.NotNil(t, err)
		require.Nil(t, result)
		assert.Contains(t, err.Error(), op.Expression)
		assert.Contains(t, err.Error(), op.Operator)
	})
}

func BenchmarkAdd(b *testing.B) {
	e := calculator.NewEngine()

	// run the Add function b.N times
	for i := 0; i < b.N; i++ {
		e.Add(2, 3)
	}
}
