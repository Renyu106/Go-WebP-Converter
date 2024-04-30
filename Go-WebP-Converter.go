package main

import (
	"bufio"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/chai2010/webp"
)

func main() {
	path := promptForPath()
	quality := promptForQuality()

	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	if fileInfo.IsDir() {
		convertDirectory(path, quality)
	} else {
		convertFile(path, quality)
	}
}

func promptForPath() string {
	reader := bufio.NewReader(os.Stdin)
	var path string

	for {
		fmt.Print("Enter file or directory path: ")
		path, _ = reader.ReadString('\n')
		path = strings.TrimSpace(path)

		if _, err := os.Stat(path); err == nil {
			break
		} else {
			fmt.Println("Invalid path. Please enter a valid file or directory path.")
			continue
		}
	}
	return path
}

func promptForQuality() int {
	reader := bufio.NewReader(os.Stdin)
	var quality int

	for {
		fmt.Print("Enter quality (1-100): ")
		qStr, _ := reader.ReadString('\n')
		qStr = strings.TrimSpace(qStr)
		var err error
		quality, err = strconv.Atoi(qStr)
		if err == nil && quality >= 1 && quality <= 100 {
			break
		} else {
			fmt.Println("Invalid quality. Please enter a number between 1 and 100.")
			continue
		}
	}
	return quality
}

func convertFile(filePath string, quality int) {
	inputFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	img, format, err := image.Decode(inputFile)
	if err != nil {
		log.Fatalf("Failed to decode image: %v", err)
	}
	log.Printf("Decoded format: %s", format)

	outputPath := filePath + ".webp"
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	options := &webp.Options{Lossless: false, Quality: float32(quality)}
	err = webp.Encode(outputFile, img, options)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Saved %s to %s\n", filePath, outputPath)
}

func convertDirectory(directoryPath string, quality int) {
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		log.Fatal(err)
	}

	var imageFiles []fs.DirEntry
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
			imageFiles = append(imageFiles, file)
		}
	}

	total := len(imageFiles)
	for i, file := range imageFiles {
		fmt.Printf("Processing file %d of %d: %s\n", i+1, total, file.Name())
		convertFile(filepath.Join(directoryPath, file.Name()), quality)
	}
}
