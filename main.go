package main

import (
	"fmt"
	"jpg-png/convert"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	jpegFolder := "./jpg"
	jpegFiles, err := readJPEGFiles(jpegFolder)
	if err != nil {
		fmt.Println("Error reading JPEG files:", err)
		return
	}

	convertJPEGFileConcurrently(jpegFiles)
}

func readJPEGFiles(jpegFolder string) (map[string][]byte, error) {
	jpegFiles := make(map[string][]byte)

	err := filepath.Walk(jpegFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (filepath.Ext(info.Name()) == ".jpg" || filepath.Ext(info.Name()) == ".jpeg") {
			jpegBytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			jpegFiles[path] = jpegBytes
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return jpegFiles, nil
}

func convertJPEGFileConcurrently(jpegFiles map[string][]byte) {
	var wg sync.WaitGroup

	progress := make(chan int, len(jpegFiles))
	go printProgressBar(progress, len(jpegFiles))

	for name, bytes := range jpegFiles {
		wg.Add(1)
		go func(name string, bytes []byte) {
			defer wg.Done()
			fmt.Println("converting JPEG File:", name)

			if err := convert.ValidateImgFileType(name, "jpg"); err != nil {
				log.Fatal(err)
			}

			convertedImg, err := convert.ConvertJpgToPng(bytes)
			if err != nil {
				log.Fatal(err)
			}

			trimmedName := strings.TrimSuffix(name, ".jpg")
			trimmedName2 := strings.TrimSuffix(trimmedName, ".jpeg")
			finalName := strings.TrimPrefix(trimmedName2, "jpg/")

			os.WriteFile(fmt.Sprintf("./png/%s.png", finalName), convertedImg, 0644)
			progress <- 1

		}(name, bytes)
	}

	go func() {
		wg.Wait()
		close(progress)
	}()

	wg.Wait()
	log.Println("Conversion was successful.Your converted files will appear in ./png folder")
}

func printProgressBar(progress <-chan int, totalTasks int) {
	fmt.Print("Progress: [")
	defer fmt.Println("]")

	completed := 0
	for range progress {
		completed++
		// Print a continuously updating progress bar
		fmt.Printf("\r[%s] %.2f%%", strings.Repeat("=", completed), float64(completed)/float64(totalTasks)*100)
	}
}
