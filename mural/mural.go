package mural

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// Start will read all pngs in input dir, perform a pixel-sort on each, then save to output dir
func Start(config Config) {
	log.Printf("Starting mural...\n")

	inputFilePaths, err := readInputDir(config.InputDir)
	if err != nil {
		log.Printf("Failed to read input directory: %s\n", err)
		os.Exit(1)
	}

	numThreads := runtime.NumCPU() // max concurrent goroutines to process images
	numImages := len(inputFilePaths)

	if numImages == 0 {
		log.Printf("Mural found no images to process\n")
		return
	}

	// avoid having excess number of goroutines
	if numThreads > numImages {
		numThreads = numImages
	}

	queue := make(chan string, numImages)
	var wg sync.WaitGroup

	// start numThreads number of goroutines
	log.Printf("Starting %d threads to process %d images\n", numThreads, numImages)
	wg.Add(numThreads)
	for i := 0; i < numThreads; i++ {
		go func() {
			// loop waiting for jobs in queue
			// exit when queue is closed
			for {
				filePath, ok := <-queue
				if !ok {
					wg.Done()
					return
				}
				processImage(filePath, config.InputDir, config.OutputDir)
			}
		}()
	}

	// add all images to queue
	for _, filePath := range inputFilePaths {
		queue <- filePath
	}

	close(queue)
	wg.Wait() // wait for all threads to finish

	log.Printf("Mural complete!\n")
}

func readInputDir(dirPath string) ([]string, error) {
	imageFilePaths := make([]string, 0)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			log.Printf("Ignoring directory: %s\n", path)
			return nil
		}

		if filepath.Ext(path) != ".png" {
			log.Printf("Ignoring non-image file: %s\n", path)
			return nil
		}

		log.Printf("Found image file: %s\n", path)
		imageFilePaths = append(imageFilePaths, path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return imageFilePaths, err
}

func processImage(filePath, inputDir, outputDir string) {
	fileName := filepath.Base(filePath)
	err := sortImage(filePath, outputDir)
	if err != nil {
		log.Printf("Error while processing %s: %s\n", fileName, err)
	}
}
