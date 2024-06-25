package controllers

import (
	"log"
	"net/http"
	"sync"

	"github.com/GoCEP/api/cep/services"
	"github.com/GoCEP/api/cep/structs"
	"github.com/GoCEP/internal/internalRouter"
	"github.com/GoCEP/internal/validator"
)

type CepController struct {
	cepService      services.CepService
	updateDataMutex sync.Mutex
	isUpdatingData  bool
}

func NewCepController(cepService services.CepService) *CepController {
	return &CepController{
		cepService: cepService,
	}
}

func (c *CepController) Read(w internalRouter.ResponseWriter, r *http.Request) {
	id := r.Context().Value(internalRouter.ContextKey("id")).(string)

	ctx := r.Context()

	data, err := c.cepService.Read(ctx, id)
	if err != nil {
		w.JSONResponse(http.StatusBadRequest, err)
		return
	}

	if data == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.JSONResponse(200, data)
}

func (c *CepController) Create(w internalRouter.ResponseWriter, r *http.Request) {
	var cep structs.Cep

	validationErrors, err := validator.ValidateBody(r, &cep, structs.CepFieldsToJSONMap)

	if validationErrors != nil {
		w.JSONResponse(http.StatusBadRequest, validationErrors)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	err = c.cepService.Create(ctx, cep)
	if err != nil {
		w.JSONResponse(400, err.Error())
		return
	}

	w.JSONResponse(http.StatusOK, cep)
}

func (c *CepController) UpdateData(w internalRouter.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	c.updateDataMutex.Lock()
	if c.isUpdatingData {
		c.updateDataMutex.Unlock()
		w.StringResponse(http.StatusConflict, "Outra atualização ja está em progresso")
		return
	}
	c.isUpdatingData = true
	c.updateDataMutex.Unlock()

	go func() {
		if err := c.cepService.UpdateData(ctx); err != nil {
			log.Printf("Error updating data: %v", err)
		}
		c.isUpdatingData = false
	}()

	w.WriteHeader(http.StatusOK)
}
