package main

import (
	"fmt"
	"sync"
)
const workerCount = 3
func greet(id int, smap *sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()
	g := fmt.Sprintf("Hello, friend! I'm Goroutine %d.", id)
	smap.Store(id, g)
}

func main() {
	var smap sync.Map
	var wg sync.WaitGroup
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go greet(i, &smap, &wg)
	}
	wg.Wait()
	smap.Range(func(key, value any) bool {
		fmt.Println(value)
		return true
	})
	fmt.Println("Goodbye, friend!")
}