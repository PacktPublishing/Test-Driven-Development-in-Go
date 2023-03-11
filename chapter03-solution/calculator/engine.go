package calculator

import (
	"fmt"

	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter03-solution/format"
)

// Operation is the wrapper object that contains
// the operator and operand of a mathematical expression.
type Operation struct {
	Expression string
	Operator   string
	Operands   []float64
}

// Engine is the mathematical logic part of the calculator.
type Engine struct {
	expectedLength  int
	validOperations map[string]func(x, y float64) (*float64, error)
}

// NewEngine initialises an engine instance and returns it.
func NewEngine() *Engine {
	e := Engine{
		expectedLength:  2,
		validOperations: make(map[string]func(x float64, y float64) (*float64, error)),
	}
	// validOperations is the map of valid operators and their corresponding functions
	e.validOperations["+"] = e.Add
	e.validOperations["-"] = e.Sub
	e.validOperations["/"] = e.Div
	e.validOperations["*"] = e.Mult
	return &e
}

// GetNumOperands returns the expected number of operands that the engine can process.
func (e *Engine) GetNumOperands() int {
	return e.expectedLength
}

// GetValidOperators returns a slice of the valid operations that the engine accepts.
func (e *Engine) GetValidOperators() []string {
	var ops []string
	for o := range e.validOperations {
		ops = append(ops, o)
	}

	return ops
}

// ProcessOperation processes a given operation and invokes the result formatter
func (e *Engine) ProcessOperation(operation Operation) (*string, error) {
	f, ok := e.validOperations[operation.Operator]
	if !ok {
		err := fmt.Errorf("no operation for operator %s found", operation.Operator)
		return nil, format.Error(operation.Expression, err)
	}
	res, err := f(operation.Operands[0], operation.Operands[1])
	if err != nil {
		return nil, format.Error(operation.Expression, err)
	}
	fres := format.Result(operation.Expression, *res)
	return &fres, nil
}

// Add is the function that processes the addition operation
func (e *Engine) Add(x, y float64) (*float64, error) {
	result := x + y
	return &result, nil
}

// Sub is the function that processes the subtract operation
func (e *Engine) Sub(x, y float64) (*float64, error) {
	result := x - y
	return &result, nil
}

// Sub is the function that processes the subtract operation
func (e *Engine) Mult(x, y float64) (*float64, error) {
	result := x * y
	return &result, nil
}

// Div is the function that processes the divide operation
func (e *Engine) Div(x, y float64) (*float64, error) {
	if y == 0 {
		return nil, fmt.Errorf("cannot divide by zero")
	}
	result := x / y
	return &result, nil
}
