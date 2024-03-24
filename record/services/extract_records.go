package services

import (
	"wannabe/record/actions"
	"wannabe/record/common"
)

func ExtractRecords(bodyBytes []byte) ([]common.Record, error) {
	return actions.ExtractRecords(bodyBytes)
}
