package services

import (
	"wannabe/config"
	"wannabe/record/actions"
	"wannabe/record/entities"
)

func ValidateRecords(wannabe config.Wannabe, records []entities.Record) ([]string, error) {
	return actions.ValidateRecords(wannabe, records)
}
