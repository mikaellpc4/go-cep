package insertData

import (
	"log"
	"os"
	"strconv"
)

func CleanJSON(unprocessedFilesChan <-chan []string, filesChan chan<- [][]byte, doneChan chan<- bool, doneZipChan <-chan bool, setCleanJSONProgress func(float64)) {
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
			for i, filePath := range filesPath {
				fileContent, err := os.ReadFile(filePath)
				if err != nil {
					log.Printf("Error reading file %s: %v", filePath, err)
					continue
				}

				cleanedData := cleanJSON(fileContent)

				batch = append(batch, cleanedData)

				os.Remove(filePath)

				percentage := float64(i + 1) / float64(len(filesPath)) 
				setCleanJSONProgress(percentage)

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
