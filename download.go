package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func download(name string) string {
	start := time.Now()
	originImageChannel := make(chan []byte)
	cdnImageChannel := make(chan []byte)

	go func() {
		originImageChannel <- getImage(
			fmt.Sprintf("%s/%s", "https://deepsleepsounds.nyc3.digitaloceanspaces.com/resources/images", name),
		)
		close(originImageChannel)
	}()
	go func() {
		cdnImageChannel <- getImage(
			fmt.Sprintf("%s/%s", "https://deepsleepsounds.nyc3.cdn.digitaloceanspaces.com/resources/images", name),
		)
		close(cdnImageChannel)
	}()

	var image []byte

	select {
	case originImage := <-originImageChannel:
		fmt.Println("Received from origin image channel:", name)
		image = originImage
	case cdnImage := <-cdnImageChannel:
		fmt.Println("Received from cdn image channel:", name)
		image = cdnImage
	}

	outputLocation := fmt.Sprintf("images/%s", name)
	permissions := 0o666
	err := os.WriteFile(outputLocation, image, os.FileMode(permissions))
	if err != nil {
		log.Fatal(err)
	}

	elapsed := time.Since(start)
	fmt.Printf("[download] %s took %s\n", name, elapsed)
	return outputLocation
}
