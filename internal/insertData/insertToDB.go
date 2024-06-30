package insertData

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/GoCEP/api/cep/repository"
	"github.com/GoCEP/api/cep/structs"
)

func InsertToDB(repo repository.CepRepositary, filesJSON chan [][]byte, doneChan <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(filesJSON)

	for {
		select {
		case files, ok := <-filesJSON:
			if !ok {
				return
			}
			var ceps []structs.Cep
			for _, fileJSON := range files {
				var cep structs.Cep
				err := json.Unmarshal(fileJSON, &cep)
				if err != nil {
					continue
				}

				ceps = append(ceps, cep)
			}
			err := repo.CreateMany(context.Background(), ceps)
			if err != nil {
				log.Fatal(err)
			}
		case done := <-doneChan:
			if done {
				return
			}
		}
	}
}
