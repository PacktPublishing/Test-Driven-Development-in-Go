package main

import (
	"fmt"
	"sync"
)
const workerCount = 3
func greet(id int, smap *sync.Map, done chan struct{}) {
	g := fmt.Sprintf("Hello, friend! I'm Goroutine %d.", id)
	smap.Store(id, g)
	done <- struct{}{}
}

func main() {
	var smap sync.Map
	done := make(chan struct{})
	for i := 0; i < workerCount; i++ {
		go greet(i, &smap, done)
	}
	for i := 0; i < workerCount; i++ {
		<-done
	}
	smap.Range(func(key, value any) bool {
		fmt.Println(value)
		return true
	})
	fmt.Println("Goodbye, friend!")
}
