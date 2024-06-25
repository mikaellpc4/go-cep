package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/GoCEP/api/cep/controllers"
	"github.com/GoCEP/api/cep/repository/implementations"
	"github.com/GoCEP/api/cep/routes"
	"github.com/GoCEP/api/cep/services"
	"github.com/GoCEP/initializers"
	"github.com/GoCEP/internal/internalRouter"
)

type MyResponse struct {
	Message string `json:"message"`
}

func main() {
	initializers.LoadEnv()

	newRouter := internalRouter.NewRouter()

	cepRepo := implementations.NewSqliteCepRepo()
	cepService := services.NewCepService(cepRepo)

	hasCepData := initializers.HasCepData()

	if !hasCepData {
		ctx := context.Background()

		err := cepService.UpdateData(ctx)
		if err != nil {
			panic(err)
		}
	}

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
