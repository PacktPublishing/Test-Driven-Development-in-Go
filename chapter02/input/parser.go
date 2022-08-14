package input

import "github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter02/calculator"

type Parser struct {
	engine *calculator.Engine
}

func (p *Parser) ProcessExpression(expr string) error {
	// implementation code
	return nil
}
