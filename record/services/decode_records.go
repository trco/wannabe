package services

import (
	"wannabe/record/actions"
	"wannabe/record/common"
)

func DecodeRecords(encodedRecords [][]byte) ([]common.Record, error) {
	return actions.DecodeRecords(encodedRecords)
}
