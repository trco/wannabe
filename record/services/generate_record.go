package services

import (
	"wannabe/config"
	"wannabe/record/actions"
	"wannabe/record/entities"
)

func GenerateRecord(config config.Records, payload entities.GenerateRecordPayload) ([]byte, error) {
	record, err := actions.GenerateRecord(config, payload)
	if err != nil {
		return nil, err
	}

	return record, nil
}
