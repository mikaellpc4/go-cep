package initializers

import (
	"os"
)

func HasCepData() bool {
	dir, _ := os.Getwd()

  cepsDir := dir + os.Getenv("CEP_DIR")

	_, err := os.Stat(cepsDir)

	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

  panic("erro ao verificar pasta de ceps: " + err.Error())
}
