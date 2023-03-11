package main

import (
	"fmt"
)

type Number interface {
	~int64 | ~float64
}

func sum[T Number](x, y T) T {
	return x + y
}

func main() {
	fmt.Println(sum[int64](2, 3))
}
