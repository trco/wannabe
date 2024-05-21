package actions

import (
	"wannabe/types"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func ValidateRecords(records []types.Record) ([]string, error) {
	var validationErrors []string

	validate = validator.New()

	for i := range records {
		err := validate.Struct(records[i])
		if err != nil {
			validationErrors = append(validationErrors, err.Error())
			continue
		}

		validationErrors = append(validationErrors, "")
	}

	return validationErrors, nil
}
