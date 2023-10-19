package main

import (
	"fmt"
	"os"
	"path/filepath"
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

	for filename := range jpegFiles {
		wg.Add(1)
		go func(filename string) {
			defer wg.Done()
			fmt.Println("JPEG File:", filename)
		}(filename)
	}

	wg.Wait()
}
