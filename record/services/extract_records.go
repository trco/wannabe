package services

import (
	"wannabe/record/actions"
	"wannabe/types"
)

func ExtractRecords(requestBody [][]byte) ([]types.Record, error) {
	return actions.ExtractRecords(requestBody)
}
