package insertData

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/mholt/archiver/v4"
)

func CountFilesInZip(zipFile string) (int, error) {
	format := archiver.Zip{}

	var fileCount int

	handler := func(ctx context.Context, f archiver.File) error {
		fileCount++
		return nil
	}

	file, err := os.Open(zipFile)
	if err != nil {
		return 0, fmt.Errorf("error reading file: %v", err)
	}
	defer file.Close()

	ctx := context.Background()

	// Extract with a nil handler to only count files
	err = format.Extract(ctx, file, nil, handler)
	if err != nil {
		return 0, fmt.Errorf("error counting files in zip: %v", err)
	}

	return fileCount, nil
}

func cleanJSON(file []byte) []byte {
	// Find the last occurrence of '}' in the byte slice
	endIndex := bytes.LastIndexByte(file, '}')

	// If '}' is found, slice the byte slice up to '}' (inclusive)
	if endIndex != -1 {
		file = file[:endIndex+1]
	}

	return file
}

func UnzipCeps(zipFile string, filesChan chan<- [][]byte, doneChan chan<- bool, setProgress func(percentage float64)) {
	defer close(doneChan)

	format := archiver.Zip{}
	var batch [][]byte

	batchSize, _ := strconv.Atoi(os.Getenv("MAX_BATCH_SIZE"))

	totalFiles, err := CountFilesInZip(zipFile)
	if err != nil {
		log.Printf("Error counting files in zip: %v", err)
		return
	}

  index := 1

	handler := func(ctx context.Context, f archiver.File) error {
		file, err := f.Open()
		if err != nil {
			log.Printf("Error reading file: %v", err)
		}
		defer file.Close()

    percentage := float64(index)/float64(totalFiles)
    setProgress(percentage)
    index++

		data, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		cleanedData := cleanJSON(data)

		batch = append(batch, cleanedData)

		if len(batch) >= batchSize {
			filesChan <- batch
			batch = nil
		}

		return nil
	}

	file, err := os.Open(zipFile)
	if err != nil {
		log.Printf("Error reading file: %v", err)
	}

	ctx := context.Background()

	err = format.Extract(ctx, file, nil, handler)
	if err != nil {
		log.Printf("Error extracting file: %v", err)
	}

	if len(batch) > 0 {
		filesChan <- batch
		batch = nil
	}

	doneChan <- true
}
