package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/GoCEP/internal/router"
)

type MyResponse struct {
	Message string `json:"message"`
}

func main() {
	newRouter := router.NewRouter()

	newRouter.POST("/cep", func(w router.ResponseWriter, r *http.Request) {
    response := MyResponse{
      Message: "DEU CERTO AEEEEE",
    }
    w.JSONResponse(200, response)
	})

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
