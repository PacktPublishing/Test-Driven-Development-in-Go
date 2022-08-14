package calculator

type Engine struct{}

func (e *Engine) Add(x, y float64) float64 {
	return x + y
}

// ... other methods
