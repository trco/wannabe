package services

import (
	"wannabe/record/actions"
	"wannabe/record/entities"
)

func ExtractRecords(bodyBytes []byte) ([]entities.Record, error) {
	return actions.ExtractRecords(bodyBytes)
}
