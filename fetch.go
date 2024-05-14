package main

import (
	"fmt"
	"sync"
	"time"
)

var files = []string{
	"SPRING_PEEPERS.jpg",
	"TWINKLE_TWINKLE_LITTLE_STAR.jpg",
	"RACE_CARS_STOCK.jpg",
	"BUBBLING_AQUARIUM.jpg",
	"OCEAN_FERRY_RIDE.jpg",
}

func downloadOne() {
	filename := "SPRING_PEEPERS.jpg"
	location := download(filename)
	resize3xAsync(location)
}

func downloadMultiple() {
	start := time.Now()

	for _, file := range files {
		location := download(file)
		resize3x(location)
	}

	elapsed := time.Since(start)
	fmt.Printf("[downloadMultiple] took %s\n", elapsed)
}

func downloadMultipleAsync() {
	start := time.Now()

	count := len(files)
	done := make(chan bool, count)
	for _, file := range files {
		go func(file string) {
			location := download(file)
			resize3xAsync(location)
			done <- true
		}(file)
	}

	for i := 0; i < count; i++ {
		<-done
	}

	close(done)

	elapsed := time.Since(start)
	fmt.Printf("[downloadMultipleAsync] took %s\n", elapsed)
}

func downloadMultipleAsyncWaitGroup() {
	start := time.Now()

	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)

		go func(file string) {
			location := download(file)
			resize3xAsync(location)

			wg.Done()
		}(file)
	}

	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("[downloadMultipleAsyncWaitGroup] took %s\n", elapsed)
}

func downloadMultipleAsyncWaitGroupSynchronized() {
	start := time.Now()

	down := func(ch chan string, file string) {
		location := download(file)
		ch <- location
	}

	res := func(ch chan string) {
		location := <-ch
		resize3xAsync(location)
	}

	ch := make(chan string)
	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)

		go func(file string) {
			defer wg.Done()

			go down(ch, file)
			res(ch)
		}(file)
	}

	wg.Wait()
	close(ch)

	elapsed := time.Since(start)
	fmt.Printf("[downloadMultipleAsyncWaitGroupSynchronized] took %s\n", elapsed)
}
