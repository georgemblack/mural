package main

import (
	"flag"
	"log"
	"os"

	"github.com/georgemblack/mural"
)

const (
	defaultInputDir  = "./data/input"
	defaultOutputDir = "./data/output"
	logFilePath      = "./mural.log"
)

func main() {
	inputDir := flag.String("input-dir", defaultInputDir, "directory containing source images")
	outputDir := flag.String("output-dir", defaultOutputDir, "directory for resulting images")
	logToFile := flag.Bool("log-to-file", false, "write log to file")

	flag.Parse()

	initLog(*logToFile)
	mural.Start(*inputDir, *outputDir)
}

func initLog(logToFile bool) {
	if !logToFile {
		log.SetOutput(os.Stdout)
		return
	}

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error initializing log file: %s\n", err)
	}
	log.SetOutput(logFile)
}
