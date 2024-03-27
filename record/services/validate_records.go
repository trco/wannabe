package services

import (
	"wannabe/config"
	"wannabe/record/actions"
	"wannabe/record/entities"
)

func ValidateRecords(config config.Config, records []entities.Record) ([]string, error) {
	return actions.ValidateRecords(config, records)
}
