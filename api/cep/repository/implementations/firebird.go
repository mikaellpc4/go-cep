package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"log"
  "os"

	"github.com/GoCEP/api/cep/structs"
	_ "github.com/nakagami/firebirdsql"
)

type FirebirdCepRepository struct {
	db *sql.DB
}

func NewFirebirdCepRepo() *FirebirdCepRepository {

	host := os.Getenv("FIREBIRD_HOST")
	port := os.Getenv("FIREBIRD_PORT")
	user := os.Getenv("FIREBIRD_USER")
	password := os.Getenv("FIREBIRD_PASSWORD")
	path := os.Getenv("FIREBIRD_PATH")

	dsn := fmt.Sprintf("%s:%s@%s:%s%s", user, password, host, port, path)
  fmt.Println(dsn)

	db, err := sql.Open("firebirdsql", dsn)
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
      CEP VARCHAR(9) PRIMARY KEY,
      LOGRADOURO VARCHAR(100),
      COMPLEMENTO VARCHAR(100),
      BAIRRO VARCHAR(100),
      LOCALIDADE VARCHAR(100),
      UF VARCHAR(2),
      IBGE VARCHAR(7)
    )
  `

	_, err = db.Exec(query)
	if err != nil {
		error := fmt.Errorf("failed to create table cep: %w", err)
		panic(error)
	}

	return &FirebirdCepRepository{
		db: db,
	}
}

func (fr *FirebirdCepRepository) Read(ctx context.Context, cep string) (*structs.Cep, error) {
	query := `
    SELECT CEP, LOGRADOURO, COMPLEMENTO, BAIRRO, LOCALIDADE, UF, IBGE
    FROM ceps
    WHERE CEP = ?
  `

	row := fr.db.QueryRowContext(ctx, query, cep)

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

		return nil, fmt.Errorf("failed to read cep: %w", err)
	}

	return &result, nil
}

func (fr *FirebirdCepRepository) Create(ctx context.Context, cep structs.Cep) error {
	query := `
		INSERT INTO ceps (CEP, LOGRADOURO, COMPLEMENTO, BAIRRO, LOCALIDADE, UF, IBGE) 
		VALUES (?, ?, ?, ?, ?, ?, ?)
  `
	_, err := fr.db.ExecContext(
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

func (fr *FirebirdCepRepository) CreateAndUpdateMany(ctx context.Context, ceps []structs.Cep) error {
	tx, err := fr.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("couldn't begin transaction")
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Printf("rollback failed: %s", rollbackErr)
			}
		}
	}()

	stmt, err := tx.PrepareContext(ctx, `
    UPDATE OR INSERT INTO ceps (CEP, LOGRADOURO, COMPLEMENTO, BAIRRO, LOCALIDADE, UF, IBGE) 
    VALUES (?, ?, ?, ?, ?, ?, ?)
  `)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %s", err)
	}
	defer stmt.Close()

	for _, cep := range ceps {
		_, err := stmt.ExecContext(ctx,
			cep.ZipCode,
			cep.PublicPlace,
			cep.Complement,
			cep.District,
			cep.Place,
			cep.Uf,
			cep.IbgeCode,
		)
		if err != nil {
			return fmt.Errorf("could not execute statement, error: %s", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction, error: %s", err)
	}

	return nil
}

func (fr *FirebirdCepRepository) Update(ctx context.Context, cep structs.Cep) error {
	query := `
		UPDATE ceps SET 
		LOGRADOURO = ?, 
		COMPLEMENTO = ?, 
		BAIRRO = ?, 
		LOCALIDADE = ?, 
		UF = ?, 
		IBGE = ? 
		WHERE CEP = ?
  `
	_, err := fr.db.ExecContext(
		ctx,
		query,
		cep.PublicPlace,
		cep.Complement,
		cep.District,
		cep.Place,
		cep.Uf,
		cep.IbgeCode,
		cep.ZipCode,
	)

	if err != nil {
		return fmt.Errorf("failed to update cep: %w", err)
	}

	return nil
}

func (fr *FirebirdCepRepository) Delete(ctx context.Context, cep string) error {
	query := `DELETE FROM ceps WHERE CEP = ?`
	_, err := fr.db.ExecContext(ctx, query, cep)
	if err != nil {
		return fmt.Errorf("failed to delete cep: %w", err)
	}
	return nil
}
