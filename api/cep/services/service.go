package services

import (
	"context"
	"fmt"
	"os"
	"path"

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

func (cepService *CepService) CreateMany(ctx context.Context, ceps []structs.Cep) error {
	return cepService.repo.CreateMany(ctx, ceps)
}

func (cepService *CepService) Update(ctx context.Context, cep structs.Cep) error {
	return cepService.repo.Update(ctx, cep)
}

func (cepService *CepService) Delete(ctx context.Context, cep string) error {
	return cepService.repo.Delete(ctx, cep)
}

func (cepService *CepService) UpdateData(ctx context.Context) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	dataLocation := dir + os.Getenv("CEP_ZIP_DIR")

	dataDir := path.Dir(dataLocation)
	_, err = os.Stat(dataDir)
	if err != nil {
		err = os.MkdirAll(dataDir, 0755)

		if err != nil {
			return fmt.Errorf("failed to mkdir, %s | error: %s", dataDir, err)
		}
	}

	err = download.File(os.Getenv("CEP_DATA_URL"), dataLocation)
	if err != nil {
		return fmt.Errorf("failed to download CEP data: %w", err)
	}

	return nil
}
