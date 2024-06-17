package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/GoCEP/api/cep/routes"
	"github.com/GoCEP/internal/internalRouter"
)

type MyResponse struct {
	Message string `json:"message"`
}

func main() {
	newRouter := internalRouter.NewRouter()

  routes.CepRoutes(newRouter)

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
