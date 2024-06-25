package initializers

import (
	"fmt"
	"os"

	"github.com/GoCEP/internal/download"
)

func VerifyAndDownloadCepData() {
	dir, _ := os.Getwd()

	dataDir := dir + "/data"
	cepsDir := dataDir + "/ceps"
	dataZip := cepsDir + "/ceps.zip"

	if _, err := os.Stat(cepsDir); err != nil {
		err := os.Mkdir(dataDir+"/ceps", 0755)
		if err != nil {
			panic("erro ao criar a pasta de ceps")
		}
	}

	if _, err := os.Stat(dataZip); err != nil {
		fmt.Println("Baixando dados de cep")
		err = download.File(os.Getenv("CEP_DATA_URL"), dataZip)
		if err != nil {
			panic("erro ao realizar o download do repositorio do OpenCEP: " + err.Error())
		}
	}
}
