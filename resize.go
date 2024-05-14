package main

import (
	"fmt"
	"log"
	"time"
)

var sizes = map[string]int{
	"1x": 6,
	"2x": 4,
	"3x": 2,
}

func resize3x(pathToFile string) {
	start := time.Now()
	filename, ext := getFilenameAndExtension(pathToFile)

	for size, multiple := range sizes {
		output := fmt.Sprintf("%s@%s.%s", filename, size, ext)
		if err := resize(pathToFile, output, multiple); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Produced ->", output)
	}

	elapsed := time.Since(start)
	fmt.Printf("[resize3x] %s took %s\n", filename, elapsed)
}

func resize3xAsync(pathToFile string) {
	start := time.Now()
	filename, ext := getFilenameAndExtension(pathToFile)

	count := len(sizes)
	done := make(chan string, count)

	for size, multiple := range sizes {
		go func(size string, multiple int) {
			output := fmt.Sprintf("%s@%s.%s", filename, size, ext)
			if err := resize(pathToFile, output, multiple); err != nil {
				log.Fatal(err)
			}
			done <- output
		}(size, multiple)
	}

	for i := 0; i < count; i++ {
		fmt.Println("Produced ->", <-done)
	}

	close(done)
	elapsed := time.Since(start)
	fmt.Printf("[resize3xAsync] %s took %s\n", filename, elapsed)
}
