package services

import (
	"wannabe/record/actions"
	"wannabe/types"
)

func GenerateRecord(config types.Records, payload types.RecordPayload) ([]byte, error) {
	record, err := actions.GenerateRecord(config, payload)
	if err != nil {
		return nil, err
	}

	return record, nil
}
