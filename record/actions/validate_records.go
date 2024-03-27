package actions

import (
	"wannabe/config"
	"wannabe/record/entities"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func ValidateRecords(config config.Config, records []entities.Record) ([]string, error) {
	var validationErrors []string

	validate = validator.New()

	validate.RegisterValidation("host_not_matching_config_server", validateHostInRecord(config))

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

// custom validation functions
func validateHostInRecord(config config.Config) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		fieldHost := fl.Parent().FieldByName(fl.StructFieldName())
		host := fieldHost.Interface().(string)

		return host == config.Server
	}
}
