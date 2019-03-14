package main

import (
	"fmt"
)

func multiplyByTwo(in <-chan int, out chan<- int) {
	fmt.Println("Initializing goroutine...")
	num := <-in
	result := num * 2
	out <- result
}

func main() {
	out := make(chan int)
	in := make(chan int)

	// Create 3 `multiplyByTwo` goroutines.
	go multiplyByTwo(in, out)
	go multiplyByTwo(in, out)
	go multiplyByTwo(in, out)

	in <- 1
	in <- 2
	in <- 3

	// Now we wait for each result to come in
	fmt.Println(<-out)
	fmt.Println(<-out)
	fmt.Println(<-out)
}