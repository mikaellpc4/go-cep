package services

import (
	"context"

	"github.com/GoCEP/api/cep/repository"
	"github.com/GoCEP/api/cep/structs"
)

type CepService struct {
	repo repository.CepRepositary
}

func NewCepService(cepRepository repository.CepRepositary) *CepService {
	return &CepService{
		repo: cepRepository,
	}
}

func (cepService *CepService) Read(ctx context.Context, cep string) (*structs.Cep, error) {
	return cepService.repo.Read(ctx, cep)
}

func (cepService *CepService) Create(ctx context.Context, cep structs.Cep) (*structs.Cep, error) {
	return cepService.repo.Create(ctx, cep)
}

func (cepService *CepService) Update(ctx context.Context, cep structs.Cep) (*structs.Cep, error) {
	return cepService.repo.Update(ctx, cep)
}

func (cepService *CepService) Delete(ctx context.Context, cepId int) error {
	return cepService.repo.Delete(ctx, cepId)
}
