package format_test

import (
	"fmt"
	"testing"

	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter03-solution/format"
	"github.com/stretchr/testify/assert"
)

func TestResult(t *testing.T) {
	// Arrange
	result := 5.55
	expr := "2+3"

	// Act
	wrappedResult := format.Result(expr, result)

	// Assert
	assert.Contains(t, wrappedResult, expr)
	assert.Contains(t, wrappedResult, fmt.Sprint(result))
}
