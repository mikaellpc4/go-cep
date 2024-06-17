package controller

import (
	"net/http"

	"github.com/GoCEP/api/cep/structs"
	"github.com/GoCEP/internal/internalRouter"
	"github.com/GoCEP/internal/validator"
)

type CepController struct {
}

func NewCepController() *CepController {
	return &CepController{}
}

func (c *CepController) Create(w internalRouter.ResponseWriter, r *http.Request) {
	var cep structs.Cep

	validationErrors, err := validator.ValidateBody(&w, r, &cep, structs.CepFieldsToJSONMap)

  if validationErrors != nil {
    w.JSONResponse(http.StatusBadRequest, validationErrors)
    return
  }

  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

	w.StringResponse(http.StatusOK, cep.Place)
}
