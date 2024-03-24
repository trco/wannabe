package actions

import (
	"wannabe/config"
	"wannabe/record/entities"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func ValidateRecords(config config.Config, records []entities.Record) ([]entities.Validation, error) {
	var validations []entities.Validation

	validate = validator.New()

	validate.RegisterValidation("host_not_matching_config_server", validateHostInRecord(config))

	for i := range records {
		err := validate.Struct(records[i])
		if err != nil {
			validations = append(validations, entities.Validation{
				Valid: false,
				Error: err.Error(),
			})

			continue
		}

		validations = append(validations, entities.Validation{
			Valid: true,
			Error: "",
		})
	}

	return validations, nil
}

// custom validation functions
func validateHostInRecord(config config.Config) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		fieldHost := fl.Parent().FieldByName(fl.StructFieldName())
		host := fieldHost.Interface().(string)

		return host == config.Server
	}
}