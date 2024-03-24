package services

import (
	"wannabe/record/actions"
	"wannabe/record/entities"
)

func DecodeRecords(encodedRecords [][]byte) ([]entities.Record, error) {
	return actions.DecodeRecords(encodedRecords)
}
