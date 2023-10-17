package main

import (
	"fmt"
	"io"
	"jpg-png/convert"
	"log"
	"os"
)

func main() {
	imgFile, err := os.Open("./jpg/A-Cat.jpg")
	if err != nil {
		log.Fatal(err)
	}

	imgBytes, err := io.ReadAll(imgFile)
	if err != nil {
		log.Fatal(err)
	}

	converted, err := convert.ConvertJpgToPng(imgBytes)
	if err != nil {
		log.Fatal(err)
	}

	outputPath := "./png/new.png"

	// Create and write the converted image to the output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	_, err = outputFile.Write(converted)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Converted %s to %s\n", "./jpg/A-Cat.jpg", outputPath)
}
