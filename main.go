package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/GoCEP/api/cep/controllers"
	"github.com/GoCEP/api/cep/repository/implementations"
	"github.com/GoCEP/api/cep/services"
	"github.com/GoCEP/api/cep/routes"
	"github.com/GoCEP/internal/internalGit"
	"github.com/GoCEP/internal/internalRouter"
)

type MyResponse struct {
	Message string `json:"message"`
}

func main() {
	dir, _ := os.Getwd()

	dataDir := dir + "/data"

	// Check if openCEP db is present
	if _, err := os.Stat(dataDir); err != nil {
		err := internalGit.GitCloneWithDepth(
			"https://github.com/SeuAliado/OpenCEP.git",
			dataDir,
			1,
		)
		if err != nil {
      panic("erro ao realizar o download do repositorio do OpenCEP")
		}
	}

	newRouter := internalRouter.NewRouter()

	cepRepo := implementations.NewSqliteCepRepo()
	cepService := services.NewCepService(cepRepo)
	cepController := controllers.NewCepController(*cepService)

	routes.CepRoutes(newRouter, *cepController)

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("error starting server: %s\n", err)
	}

	fmt.Println("🐱💻 GoCEP server started on", l.Addr().String())

	if err := http.Serve(l, newRouter); err != nil {
		fmt.Printf("server closed: %s\n", err)
	}

	os.Exit(1)
}
