package main

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/image/draw"
)

func createDirectory() {
	paths := []string{"images", "images/thumbnails"}

	for _, path := range paths {
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(path, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func getImage(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	return bytes
}

func getFilenameAndExtension(pathToFile string) (string, string) {
	filenameWithExt := strings.Split(pathToFile, "/")[1]
	filenameAndExt := strings.Split(filenameWithExt, ".")

	filename := filenameAndExt[0]
	extension := filenameAndExt[1]

	return filename, extension
}

func resize(pathToFile, pathToOutputImage string, multiple int) error {
	// Open up the damn file
	input, err := os.Open(pathToFile)
	if err != nil {
		return err
	}

	// don't forget to defer closing it to avoid memory leak
	defer input.Close()

	// create a new image buffer
	imageBuff := make([]byte, 512)

	// read the original file into the image buffer, so we can check what type of file it is
	_, err = input.Read(imageBuff)
	if err != nil {
		return err
	}

	// get the file type from the buffer
	fileType := http.DetectContentType(imageBuff)

	// create the output file
	output, _ := os.Create(fmt.Sprintf("images/thumbnails/%s", pathToOutputImage))

	// defer closing the output file until we are done writing it and the function exits (avoid memory leak)
	defer output.Close()

	// seek the file back to the beginning or else we won't be able to write the whole file
	_, err = input.Seek(0, 0)
	if err != nil {
		return err
	}

	// create a new image variable
	var src image.Image

	// determine if the original file was a png or jpeg before continuing
	// to decode the image (from PNG/JPEG to image.Image):
	if strings.EqualFold(fileType, "image/png") {
		src, err = png.Decode(input)
		if err != nil {
			return err
		}
	} else {
		src, err = jpeg.Decode(input)
		if err != nil {
			return err
		}
	}

	// create a whole new sized image
	// Set the expected size that you want:
	destinationImage := image.NewRGBA(image.Rect(0, 0, src.Bounds().Max.X/multiple, src.Bounds().Max.Y/multiple))

	// Resize:
	draw.NearestNeighbor.Scale(destinationImage, destinationImage.Rect, src, src.Bounds(), draw.Over, nil)

	// Encode to `output`:
	if strings.EqualFold(fileType, "image/png") {
		if err := png.Encode(output, destinationImage); err != nil {
			return err
		}
	} else {
		if err := jpeg.Encode(output, destinationImage, nil); err != nil {
			return err
		}
	}

	return nil
}
