package main

import (
	"fmt"
	"sync"
)

func safelyClose(once *sync.Once, ch chan struct{}) {
	fmt.Println("Hello, friend!")
	once.Do(func() {
		fmt.Println("Channel closed.")
		close(ch)
	})
}

func main() {
	var once sync.Once
	ch := make(chan struct{})
	for i := 0; i<3; i++ {
		go safelyClose(&once, ch)
	}
	<-ch
	fmt.Println("Goodbye, friend!")
}
