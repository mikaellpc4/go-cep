package services

import (
	"context"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/GoCEP/api/cep/repository"
	// "github.com/GoCEP/internal/download"
	"github.com/GoCEP/internal/insertData"
)

type CepService struct {
	repos []repository.CepRepositary
}

func NewCepService(cepRepository []repository.CepRepositary) *CepService {
	return &CepService{
		repos: cepRepository,
	}
}

func (cepService *CepService) UpdateData(ctx context.Context) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	dataLocation := dir + os.Getenv("CEP_ZIP_LOCATION")

	dataDir := path.Dir(dataLocation)
	_, err = os.Stat(dataDir)
	if err != nil {
		err = os.MkdirAll(dataDir, 0755)

		if err != nil {
			return fmt.Errorf("failed to mkdir, %s | error: %s", dataDir, err)
		}
	}

	// err = download.File(os.Getenv("CEP_DATA_URL"), dataLocation)
	// if err != nil {
	// 	return fmt.Errorf("failed to download CEP data: %w", err)
	// }

	unprocessedFilesChan := make(chan []string)
	filesJSON := make(chan [][]byte)
	doneZipChan := make(chan bool)
	doneChan := make(chan bool)
	var wg sync.WaitGroup

	wg.Add(1)
	go insertData.UnzipCeps(dataLocation, unprocessedFilesChan, doneZipChan, &wg)

	wg.Add(1)
	go insertData.CleanJSON(unprocessedFilesChan, filesJSON, doneChan, doneZipChan, &wg)

	for _, repo := range cepService.repos {
		wg.Add(1)
		go insertData.InsertToDB(repo, filesJSON, doneChan, &wg)
	}

	wg.Wait()
	close(filesJSON)

	fmt.Println("finished CEP Update")

	return nil
}
