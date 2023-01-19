package main

import (
	"fmt"
	"sync"
)

var greetings []string

const workerCount = 3

func greet(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	g := fmt.Sprintf("Hello, friend! I'm Goroutine %d.", id)
	greetings = append(greetings, g)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go greet(i, &wg)
	}
	wg.Wait()
	for _, g := range greetings {
		fmt.Println(g)
	}
	fmt.Println("Goodbye, friend!")
}
