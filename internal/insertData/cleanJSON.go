package insertData

import (
	"bytes"
	"log"
	"os"
	"strconv"
	"sync"
)

func cleanJSON(file []byte) []byte {
	// Find the last occurrence of '}' in the byte slice
	endIndex := bytes.LastIndexByte(file, '}')

	// If '}' is found, slice the byte slice up to '}' (inclusive)
	if endIndex != -1 {
		file = file[:endIndex+1]
	}

	return file
}

func CleanJSON(unprocessedFilesChan <-chan []string, filesChan chan<- [][]byte, doneChan chan<- bool, doneZipChan <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(doneChan)

	batchEnv := os.Getenv("MAX_BATCH_SIZE")
	batchSize, err := strconv.Atoi(batchEnv)
	if err != nil {
		panic("Probably invalid MAX_BATCH_SIZE in .env")
	}
	var batch [][]byte

	stop := false

	for !stop {
		select {
		case filesPath, ok := <-unprocessedFilesChan:
			if !ok {
				return
			}
			for _, filePath := range filesPath {
				fileContent, err := os.ReadFile(filePath)
				if err != nil {
					log.Printf("Error reading file %s: %v", filePath, err)
					continue
				}

				cleanedData := cleanJSON(fileContent)

				batch = append(batch, cleanedData)

				os.Remove(filePath)

				if len(batch) >= batchSize {
					filesChan <- batch
					batch = nil
				}
			}
		case done := <-doneZipChan:
			if done {
				stop = true
			}
		}
	}

	if len(batch) > 0 {
		filesChan <- batch
	}

	doneChan <- true
}
