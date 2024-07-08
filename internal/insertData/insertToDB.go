package insertData

import (
	"context"
	"encoding/json"
	"log"

	"github.com/GoCEP/api/cep/repository"
	"github.com/GoCEP/api/cep/structs"
)

func InsertToDB(repo repository.CepRepositary, filesJSON chan [][]byte, doneChan <-chan bool, setProgress func(float64)) {
	for {
		select {
		case files, ok := <-filesJSON:
			if !ok {
				return
			}
			var ceps []structs.Cep
			for i, fileJSON := range files {
				var cep structs.Cep
				err := json.Unmarshal(fileJSON, &cep)
				if err != nil {
					continue
				}
				percentage := float64(i+1) / float64(len(files))
				setProgress(percentage)

				ceps = append(ceps, cep)
			}

			err := repo.CreateAndUpdateMany(context.Background(), ceps)
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
