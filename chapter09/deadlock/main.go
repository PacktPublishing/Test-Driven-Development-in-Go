package main

import (
	"fmt"
	"sync"
)

var greetings []string

const workerCount = 3

func greet(id int, ch chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	g := fmt.Sprintf("Hello, friend! I'm Goroutine %d.", id)
	<-ch
	greetings = append(greetings, g)
	ch <- struct{}{}
}

func main() {
	ch := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go greet(i, ch, &wg)
	}
	ch <- struct{}{}
	wg.Wait()
	for _, g := range greetings {
		fmt.Println(g)
	}
	fmt.Println("Goodbye, friend!")
}
