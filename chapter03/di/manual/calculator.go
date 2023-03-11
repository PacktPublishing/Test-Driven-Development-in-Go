package main

import "fmt"

type Adder interface {
	Add(x, y float64) float64
}

type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

func (e Engine) Add(x, y float64) float64 {
	return x + y
}

type Calculator struct {
	Adder Adder
}

func NewCalculator(a Adder) *Calculator {
	return &Calculator{Adder: a}
}

func (c Calculator) PrintAdd(x, y float64) {
	fmt.Println("Result:", c.Adder.Add(x, y))
}

func main() {
	engine := NewEngine()
	calc := NewCalculator(engine)

	calc.PrintAdd(2.5, 6.3)
}
