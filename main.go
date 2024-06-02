package main

import (
	"bytes"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var files = []string{}

func main() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU) // Utilize all available CPU cores

	filepath.Walk("assets/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	// Measure the execution time for the parallel implementation
	start := time.Now()

	extractChan := extract(files, 8)
	convertedChan := convert(extractChan, 8)
	compressedChan := compress(convertedChan, 8)
	save(compressedChan)

	parallelDuration := time.Since(start)

	// Measure the execution time for the sequential implementation
	start = time.Now()
	for _, image := range files {
		buf, _ := readFile(image)
		file := File{
			name:   image,
			buffer: buf,
		}
		buf, _ = file.convertToJPG()
		buf, _ = file.compress()

		outFile := strings.Replace(image, "assets", "out", 1)
		img, _ := os.Create(outFile)
		processedImage, _ := jpeg.Decode(bytes.NewReader(buf))
		jpeg.Encode(img, processedImage, &jpeg.Options{})

		log.Println("[SEQ] completed processing: ", image)
	}

	sequentialDuration := time.Since(start)

	log.Printf("Parallel implementation took: %v\n", parallelDuration)
	log.Printf("Sequential implementation took: %v\n", sequentialDuration)
}
