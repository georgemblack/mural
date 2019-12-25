package mural

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type pixel struct {
	color  color.Color
	weight float64
}

func sortImage(inputFilePath, outputDirPath string, strength int) error {
	inputFileName := filepath.Base(inputFilePath)

	// open file
	log.Printf("Opening file at: %s\n", inputFilePath)
	srcFile, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// decode
	srcImage, err := png.Decode(srcFile)
	if err != nil {
		return err
	}

	// find bounds
	bounds := srcImage.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y
	log.Printf("Image width: " + strconv.Itoa(width) + " height: " + strconv.Itoa(height) + "\n")

	// begin pixel-sort
	log.Printf("Starting sort on image: %s\n", inputFileName)
	destImage := sortImageHelper(srcImage, width, height, strength)

	// create or open file in output dir
	outputFilePath := outputDirPath + "/" + inputFileName
	destFile, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer destFile.Close()

	png.Encode(destFile, destImage)
	log.Printf("Finished sort for image: " + inputFileName)

	return nil
}

func sortImageHelper(srcImage image.Image, width, height, strength int) image.Image {
	destImage := image.NewRGBA(image.Rect(0, 0, width, height))

	currRow := make([]color.Color, width)
	sortedRow := make([]color.Color, width)

	// for each row
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			currRow[x] = srcImage.At(x, y)
		}
		sortedRow = sortRow(currRow, strength)
		for x := 0; x < width; x++ {
			destImage.Set(x, y, sortedRow[x])
		}
	}
	return destImage
}

func sortRow(row []color.Color, strength int) []color.Color {
	length := len(row)
	pixels := make([]pixel, length)
	var red, green, blue uint32
	var brightness, weight float64

	// calculate pixel weight, store in new slice
	for i := 0; i < length; i++ {
		red, green, blue, _ = row[i].RGBA()
		brightness = 0.2126*float64(red) + 0.7152*float64(green) + 0.0722*float64(blue)
		weight = (brightness / float64(65535)) + (float64(i)/float64(length-1))*float64(strength)
		pixels[i] = pixel{row[i], weight}
	}

	// sort based on pixel's weight
	sort.SliceStable(pixels, func(i, j int) bool {
		return pixels[i].weight > pixels[j].weight
	})

	// place pixels back into row
	for i := 0; i < length; i++ {
		row[i] = pixels[i].color
	}

	return row
}
