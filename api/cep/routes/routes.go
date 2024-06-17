package cep

import (
	"net/http"

	"github.com/GoCEP/internal/internalRouter"
)

func Routes(router *internalRouter.Router) {
	router.GROUP("/api", func(router *internalRouter.Router) {
		router.GET("/test", func(w internalRouter.ResponseWriter, request *http.Request) {
			responseObj := map[string]string{
				"message": "test get",
			}
			w.JSONResponse(200, responseObj)
		})
    
		router.POST("/test", func(w internalRouter.ResponseWriter, request *http.Request) {
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
