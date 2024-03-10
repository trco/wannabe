package services

import (
	"wannabe/record/actions"
	"wannabe/record/entities"
)

func ExtractRecords(body []byte) ([]entities.Record, error) {
	return actions.ExtractRecords(body)
}
