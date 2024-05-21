package services

import (
	"wannabe/record/actions"
	"wannabe/types"
)

func ValidateRecords(records []types.Record) ([]string, error) {
	return actions.ValidateRecords(records)
}
