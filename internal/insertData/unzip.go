package insertData

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"sync"

	"github.com/GoCEP/internal/progressBar"
)

func UnzipCeps(zipFile string, unprocessedFilesChan chan<- []string, doneChan chan<- bool, wg *sync.WaitGroup) {
	defer close(doneChan)
	defer wg.Done()

	r, err := zip.OpenReader(zipFile)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := r.Close(); err != nil {
			log.Printf("Failed to unzip file, error: %s", err)
		}
	}()

	if err != nil {
		log.Fatal(err)
	}

  bar := progressBar.Create(len(r.File), "Extraindo ZIP e Inserindo no banco de dados")

	var files []string
	for _, file := range r.File {
		tmpFile, err := os.CreateTemp("", "go-cep-*.tmp")
		if err != nil {
			log.Fatal(err)
		}

		rc, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(tmpFile, rc)
		if err != nil {
			log.Fatal(err)
		}
		rc.Close()
		tmpFile.Close()

		files = append(files, tmpFile.Name())
    bar.Add(1)

		if len(files) > 10000 {
			unprocessedFilesChan <- files
			files = nil
		}
	}

	if len(files) > 0 {
		unprocessedFilesChan <- files
	}

	doneChan <- true
}
