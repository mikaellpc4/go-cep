package repository

import (
	"context"

	"github.com/GoCEP/api/cep/structs"
)

type CepRepositary interface {
	Read(ctx context.Context, cep string) (*structs.Cep, error)
	Create(ctx context.Context, cep structs.Cep) error
	CreateAndUpdateMany(ctx context.Context, ceps []structs.Cep) error
	Update(ctx context.Context, cep structs.Cep) error
	Delete(ctx context.Context, cep string) error
}
