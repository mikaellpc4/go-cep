package services

import (
	"context"
	"fmt"
	"os"

	"github.com/GoCEP/api/cep/repository"
	"github.com/GoCEP/api/cep/structs"
	"github.com/GoCEP/internal/download"
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

func (cepService *CepService) Create(ctx context.Context, cep structs.Cep) error {
	return cepService.repo.Create(ctx, cep)
}

func (cepService *CepService) Update(ctx context.Context, cep structs.Cep) error {
	return cepService.repo.Update(ctx, cep)
}

func (cepService *CepService) Delete(ctx context.Context, cep string) error {
	return cepService.repo.Delete(ctx, cep)
}

func (cepService *CepService) UpdateData(ctx context.Context) error {
	fmt.Println("Baixando dados de cep")
	dir, _ := os.Getwd()

	dataDir := dir + os.Getenv("CEP_DIR")

	err := download.File(os.Getenv("CEP_DATA_URL"), dataDir)
	if err != nil {
		return fmt.Errorf("failed to download CEP data: %w", err)
	}

	return nil
}
