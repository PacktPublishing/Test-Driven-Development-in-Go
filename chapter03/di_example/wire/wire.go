//go:build wireinject

package main

import "github.com/google/wire"

func InitCalc() Calculator {
	wire.Build(NewCalculator, NewEngine)
	return Calculator{}
}
