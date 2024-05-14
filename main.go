package main

import (
	"fmt"
	"sync"
	"time"
)

func printOdd(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 20; i += 2 {
		fmt.Println(i)
		ch <- 1 // Signal to the even goroutine that it's its turn to print
		<-ch    // Wait for signal from the even goroutine before printing again
	}
}

func printEven(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 2; i <= 20; i += 2 {
		<-ch // Wait for signal from the odd goroutine before printing
		fmt.Println(i)
		ch <- 1 // Signal to the odd goroutine that it's its turn to print
	}
}

func synchronize() {
	ch := make(chan int)
	var wg sync.WaitGroup
	wg.Add(2)

	go printOdd(ch, &wg)
	go printEven(ch, &wg)

	wg.Wait()
	close(ch)
}

func firstDone() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	// Send values to channels
	go func() {
		fmt.Println("Doing work on ch1")
		time.Sleep(time.Second * 5)
		ch1 <- 1
	}()
	go func() {
		fmt.Println("Doing work on ch2")
		time.Sleep(time.Second * 2)
		ch2 <- 2
	}()

	// Wait for and print values received from channels
	select {
	case x := <-ch1:
		fmt.Println("Received from ch1:", x)
	case y := <-ch2:
		fmt.Println("Received from ch2:", y)
	}

	close(ch1)
	close(ch2)
}

func display(message string) {
	fmt.Println(message)
}

func main() {
	go display("Process 1")
	display("Process 2")
	time.Sleep(time.Second * 1)
}
