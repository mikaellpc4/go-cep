package validator

import (
	"encoding/json"
	"net/http"

	"github.com/GoCEP/internal/internalRouter"
	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field string
	Err   string
}

type ValidationErrors struct {
	Field string
	Err   string
}

func ValidateBody(
	w *internalRouter.ResponseWriter,
	r *http.Request,
	validation interface{},
	mappedValidation map[string]string,
) ([]ValidationError, error) {
	err := json.NewDecoder(r.Body).Decode(validation)

	if err != nil {
		return nil, err
	}

	validate := validator.New()
	err = validate.Struct(validation)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var errors []ValidationError

		for _, err := range validationErrors {
			jsonFieldName, ok := mappedValidation[err.StructField()]
			if !ok {
				jsonFieldName = err.StructField()
			}

			fieldError := ValidationError{
				Field: jsonFieldName,
				Err:   err.ActualTag(),
			}
			errors = append(errors, fieldError)
		}

		return errors, nil
	}

	return nil, nil
}
