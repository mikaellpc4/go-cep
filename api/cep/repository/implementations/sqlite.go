package implementations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/GoCEP/api/cep/structs"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteCepRepository struct {
	db *sql.DB
}

func NewSqliteCepRepo() *SqliteCepRepository {
	db, err := sql.Open("sqlite3", "./data/cep.sqlite")

	if err != nil {
		error := fmt.Errorf("failed to open database: %w", err)
		panic(error)
	}

	if err := db.Ping(); err != nil {
		error := fmt.Errorf("failed to ping database: %w", err)
		panic(error)
	}

	query := `
    CREATE TABLE IF NOT EXISTS ceps (
      CEP TEXT,
      LOGRADOURO TEXT,
      COMPLEMENTO TEXT,
      BAIRRO TEXT,
      LOCALIDADE TEXT,
      UF TEXT,
      IBGE TEXT
    )
  `

	_, err = db.Exec(
		query,
	)

	if err != nil {
    error := fmt.Errorf("failed to create table cep: %w", err)
    panic(error)
	}

	return &SqliteCepRepository{
		db: db,
	}
}

func (sr *SqliteCepRepository) Read(ctx context.Context, cep string) (*structs.Cep, error) {
	query := `
    SELECT CEP, LOGRADOURO, COMPLEMENTO, BAIRRO, LOCALIDADE, UF, IBGE
    FROM ceps
    WHERE CEP = ?
  `

	row := sr.db.QueryRowContext(
		ctx,
		query,
		cep,
	)

	var result structs.Cep
	err := row.Scan(
		&result.ZipCode,
		&result.PublicPlace,
		&result.Complement,
		&result.District,
		&result.Place,
		&result.Uf,
		&result.IbgeCode,
	)
  if err != nil {
    if err == sql.ErrNoRows {
      return nil, nil
    }

    return nil, fmt.Errorf("faied to read cep: %w", err)
  }

	return &result, nil
}

func (sr *SqliteCepRepository) Create(ctx context.Context, cep structs.Cep) error {
	query := `
		INSERT INTO ceps (CEP, LOGRADOURO, COMPLEMENTO, BAIRRO, LOCALIDADE, UF, IBGE) 
		VALUES (?, ?, ?, ?, ?, ?, ?)
  `
	_, err := sr.db.ExecContext(
		ctx,
		query,
		cep.ZipCode,
		cep.PublicPlace,
		cep.Complement,
		cep.District,
		cep.Place,
		cep.Uf,
		cep.IbgeCode,
	)

	if err != nil {
		return fmt.Errorf("failed to create cep: %w", err)
	}

	return nil
}

func (sr *SqliteCepRepository) Update(ctx context.Context, cep structs.Cep) error {
	return nil
}

func (sr *SqliteCepRepository) Delete(ctx context.Context, cep string) error {
	return nil
}