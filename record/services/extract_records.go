package services

import (
	"wannabe/record/actions"
	"wannabe/record/entities"
)

func ExtractRecords(requestBody [][]byte) ([]entities.Record, error) {
	return actions.ExtractRecords(requestBody)
}
