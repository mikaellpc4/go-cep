package insertData

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/GoCEP/api/cep/repository"
	"github.com/GoCEP/api/cep/structs"
)

func InsertToDB(repo repository.CepRepositary, filesChan <-chan []string, doneChan <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case batch, ok := <-filesChan:
			if !ok {
				return
			}

			for _, file := range batch {
				data, err := os.ReadFile(file)
				if err != nil {
					log.Fatal(err)
				}

				var cep structs.Cep
				err = json.Unmarshal(data, &cep)
				if err != nil {
					log.Fatal(err)
				}

				// Insert data using repository method
				err = repo.Create(context.Background(), cep)
				if err != nil {
					log.Fatal(err)
				}

				// Remove the temporary file
				err = os.Remove(file)
				if err != nil {
					log.Fatal(err)
				}
			}

		case <-doneChan:
      fmt.Println("true")
			return
		}
	}
}
