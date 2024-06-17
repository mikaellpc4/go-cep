package implementations

import (
	"context"
	"database/sql"

	"github.com/GoCEP/api/cep/structs"
)

type SqliteCepRepository struct {
	db *sql.DB
}

func NewSqliteCepRepo(db *sql.DB) *SqliteCepRepository {
	return &SqliteCepRepository{
		db: db,
	}
}

func (sr *SqliteCepRepository) Read(ctx context.Context, cep string) (*structs.Cep, error) {
	return nil, nil
}

func (sr *SqliteCepRepository) Create(ctx context.Context, cep structs.Cep) (*structs.Cep, error) {
	return nil, nil
}

func (sr *SqliteCepRepository) Update(ctx context.Context, cep structs.Cep) (*structs.Cep, error) {
	return nil, nil
}

func (sr *SqliteCepRepository) Delete(ctx context.Context, cep string) error {
	return nil
}
