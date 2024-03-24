package services

import (
	"wannabe/config"
	"wannabe/record/actions"
	"wannabe/record/common"
)

func ValidateRecords(config config.Config, records []common.Record) ([]common.Validation, error) {
	return actions.ValidateRecords(config, records)
}
