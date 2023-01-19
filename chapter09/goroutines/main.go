package main

import (
	"fmt"
)

func greet(ch chan struct{}) {
	fmt.Println("Hello, friend!")
	close(ch)
}

func main() {
	ch := make(chan struct{})
	go greet(ch)
	<-ch
	fmt.Println("Child goroutine finished.")
	fmt.Println("Goodbye, friend!")
}
