package main

import (
	"bytes"
	"image/jpeg"
	"log"
	"os"
	"strings"
	"sync"
)

const maxPoolSize = 8

func extract(files []string, numWorkers int) <-chan File {
	extractChan := make(chan File)
	go func() {
		defer close(extractChan)
		wg := sync.WaitGroup{}
		fileChan := make(chan string, numWorkers)

		for range maxPoolSize {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for image := range fileChan {
					file, _ := readFile(image)
					extractChan <- File{
						name:   image,
						buffer: file,
					}
				}
			}()
		}

		for _, image := range files {
			fileChan <- image
		}

		close(fileChan)
		wg.Wait()
	}()
	return extractChan
}

func convert(data <-chan File, numWorkers int) <-chan File {
	convertedChan := make(chan File)

	go func() {
		defer close(convertedChan)

		wg := sync.WaitGroup{}
		fileChan := make(chan File, numWorkers)

		for range maxPoolSize {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for image := range fileChan {
					file, _ := image.convertToJPG()
					convertedChan <- File{
						name:   image.name,
						buffer: file,
					}
				}
			}()
		}

		for image := range data {
			fileChan <- image
		}

		close(fileChan)
		wg.Wait()
	}()

	return convertedChan
}

func compress(data <-chan File, numWorkers int) <-chan File {
	compressChan := make(chan File)

	go func() {
		defer close(compressChan)
		wg := sync.WaitGroup{}
		fileChan := make(chan File, numWorkers)

		for range maxPoolSize {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for image := range fileChan {
					file, _ := image.compress()
					compressChan <- File{
						name:   image.name,
						buffer: file,
					}
				}
			}()
		}

		for image := range data {
			fileChan <- image
		}

		close(fileChan)
		wg.Wait()
	}()

	return compressChan
}

func save(data <-chan File) {
	var wg sync.WaitGroup

	for range maxPoolSize {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for image := range data {
				log.Println("completed processing: ", image.name)
				outFile := strings.Replace(image.name, "assets", "out", 1)
				img, _ := os.Create(outFile)
				processedImage, _ := jpeg.Decode(bytes.NewReader(image.buffer))
				jpeg.Encode(img, processedImage, &jpeg.Options{})
			}
		}()
	}

	wg.Wait()
}
