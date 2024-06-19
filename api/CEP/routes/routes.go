package cep

import (
	"net/http"

	"github.com/GoCEP/internal/internalGit"
	"github.com/GoCEP/internal/internalRouter"
)

func Routes(router *internalRouter.Router) {
	router.GROUP("/api", func(router *internalRouter.Router) {
		router.GROUP("/cep", func(router *internalRouter.Router) {
			router.GET("/:cep", func(w internalRouter.ResponseWriter, request *http.Request) {
        _ = internalGit.GitLog(".")
				cep, ok := request.Context().Value(internalRouter.ContextKey("cep")).(string)

				if !ok {
					w.StringResponse(http.StatusBadRequest, "cep value is not a string")
					return
				}

				w.StringResponse(http.StatusOK, cep)
			})
		})

		router.POST("/updateDatabase", func(w internalRouter.ResponseWriter, request *http.Request) {
			responseObj := map[string]string{
				"message": "test get",
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
