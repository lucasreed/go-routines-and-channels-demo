package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func watcher(ch1 chan string, ch2 chan string, stop chan bool) {
	// defer is a special keyworkd
	// anything that defer references gets run just before the enclosing function returns
	// this is where we put things we want to gracefully shut down
	defer func() {
		fmt.Println("Closing the channels.")
		close(ch1)
		close(ch2)
		close(stop)
		fmt.Println("Channels closed")
	}()

	// an empty for loop will run forever unless a `return` is called
	for {
		// using a select statement allows us to read from multiple channels in one loop/function
		select {
		case data1 := <-ch1:
			fmt.Println("Received from ch1")
			fmt.Println(data1)
		case data2 := <-ch2:
			fmt.Println("Received from ch2")
			fmt.Println(data2)
		case <-stop: // if we read from the stop channel, we just return which ends the for loop
			fmt.Println("Was told to stop!")
			return
		}
	}
}

func writeToChannels(ch1 chan string, ch2 chan string) {
	for {
		fmt.Println("Sending to ch1")
		ch1 <- time.Now().Format("2006-01-02 15:04:05")
		time.Sleep(3 * time.Second)

		fmt.Println("Sending to ch2")
		ch2 <- time.Now().Format("2006-01-02 15:04:05")
		time.Sleep(3 * time.Second)
	}
}

func main() {
	ch1 := make(chan string) // an unbuffered channel
	ch2 := make(chan string)
	stop := make(chan bool, 1) // buffered channel - mostly unecessary, but to show one written out

	go watcher(ch1, ch2, stop)   // go routine that watches all the channels for data
	go writeToChannels(ch1, ch2) // where the data sending gets performed

	// catch interrupt signal and gracefully shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt) // specifically look for SIGINT (ctrl-C)
	<-signalChan
	fmt.Println("Received an interrupt, stopping services...")
	stop <- true

	time.Sleep(1 * time.Second) // give the defers time to close channels
	fmt.Println("Stopped")
}
