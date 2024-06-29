package insertData

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"strconv"
)

func UnzipCeps(zipFile string, filesChan chan<- []string, doneChan chan<- bool) {
	defer close(filesChan)
	defer close(doneChan)

	r, err := zip.OpenReader(zipFile)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := r.Close(); err != nil {
			log.Printf("Failed to unzip file, error: %s", err)
		}
	}()

	batchEnv := os.Getenv("MAX_BATCH_SIZE")
	batchSize, err := strconv.Atoi(batchEnv)
	if err != nil {
		log.Fatal(err)
	}

	var batch []string
	for _, file := range r.File {
		tmpFile, err := os.CreateTemp("", "go-cep-*.tmp")
		if err != nil {
			log.Fatal(err)
		}
		defer tmpFile.Close()

		rc, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(tmpFile, rc)
		if err != nil {
			log.Fatal(err)
		}
		rc.Close()

		batch := append(batch, tmpFile.Name())
		if len(batch) >= batchSize {
			filesChan <- batch
			batch = nil
		}
	}

	if len(batch) > 0 {
		filesChan <- batch
	}

	doneChan <- true
}
