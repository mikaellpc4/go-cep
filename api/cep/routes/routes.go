package routes

import (
	"net/http"

	controller "github.com/GoCEP/api/cep/controllers"
	"github.com/GoCEP/internal/internalRouter"
)

func CepRoutes(router *internalRouter.Router) {
	cepController := controller.NewCepController()

	router.GROUP("/cep", func(router *internalRouter.Router) {
		router.POST("/", cepController.Create)

		router.POST("/", func(w internalRouter.ResponseWriter, request *http.Request) {
			responseObj := map[string]string{
				"message": "test post",
			}
			w.JSONResponse(200, responseObj)
		})
	})

	router.GET("/hello", func(w internalRouter.ResponseWriter, r *http.Request) {
		responseObj := map[string]string{
			"message": "hello",
		}
		w.JSONResponse(200, responseObj)
	})
}
