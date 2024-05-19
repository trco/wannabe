package services

import (
	"wannabe/record/actions"
	"wannabe/types"
)

func DecodeRecords(encodedRecords [][]byte) ([]types.Record, error) {
	return actions.DecodeRecords(encodedRecords)
}
