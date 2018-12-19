package mural

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// Start will read all pngs in input dir, perform a pixel-sort on each, then save to output dir
func Start(inputDir, outputDir string) {
	log.Printf("starting mural...\n")

	inputFilePaths, err := readInputDir(inputDir)
	if err != nil {
		log.Printf("failed to read input directory: %s\n", err)
		os.Exit(1)
	}

	numThreads := runtime.NumCPU() // max concurrent goroutines to process images
	numImages := len(inputFilePaths)

	if numImages == 0 {
		log.Printf("mural found no images to process\n")
		return
	}

	// avoid having excess number of goroutines
	if numThreads > numImages {
		numThreads = numImages
	}

	queue := make(chan string, numImages)
	var wg sync.WaitGroup

	// start numThreads number of goroutines
	log.Printf("starting %d threads to process %d images\n", numThreads, numImages)
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
				processImage(filePath, inputDir, outputDir)
			}
		}()
	}

	// add all images to queue
	for _, filePath := range inputFilePaths {
		queue <- filePath
	}

	close(queue)
	wg.Wait() // wait for all threads to finish

	log.Printf("mural complete!\n")
}

func readInputDir(dirPath string) ([]string, error) {
	imageFilePaths := make([]string, 0)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			log.Printf("ignoring directory: %s\n", path)
			return nil
		}

		if filepath.Ext(path) != ".png" {
			log.Printf("ignoring non-image file: %s\n", path)
			return nil
		}

		log.Printf("found image file: %s\n", path)
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
		log.Printf("error while processing %s: %s\n", fileName, err)
	}
}
