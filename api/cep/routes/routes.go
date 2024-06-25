package routes

import (
	"github.com/GoCEP/api/cep/controllers"
	"github.com/GoCEP/internal/internalRouter"
)

func CepRoutes(router *internalRouter.Router, cepController *controllers.CepController) {
	router.GROUP("/cep", func(router *internalRouter.Router) {
		router.POST("/", cepController.Create)
		router.GET("/:id", cepController.Read)
		router.POST("/update", cepController.UpdateData)
	})
}
