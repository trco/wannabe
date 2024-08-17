package actions

import (
	"github.com/trco/wannabe/types"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func ValidateRecords(records []types.Record) ([]string, error) {
	var validationErrors []string

	validate = validator.New()
	validate.RegisterValidation("content_type_header_present", validateHeaders)

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

func validateHeaders(fl validator.FieldLevel) bool {
	headers, ok := fl.Field().Interface().(map[string][]string)
	if !ok {
		return false
	}

	_, ok = headers["Content-Type"]
	return ok
}
