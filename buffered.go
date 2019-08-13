package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 2)

	go func(ch chan int) {
		for i := 1; i <= 5; i++ {
			ch <- i
			fmt.Println("Func goroutine sends data: ", i)
		}
		close(ch)
	}(ch)

	fmt.Println("Main goroutine sleeps 2 seconds")
	time.Sleep(time.Second * 2)

	fmt.Println("Main goroutine begins receiving data")
	for d := range ch {
		fmt.Println("Main goroutine received data:", d)
	}
}
