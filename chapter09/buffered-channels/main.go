package main

import (
	"fmt"
)

const workerCount = 3

func greet(id int, ch chan string) {
	g := fmt.Sprintf("Hello, friend! I'm Goroutine %d.", id)
	ch <- g
	fmt.Printf("Goroutine %d completed.\n", id)
}

func main() {
	ch := make(chan string, workerCount)
	for i := 0; i < workerCount; i++ {
		go greet(i, ch)
	}
	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println("Goodbye, friend!")
}
