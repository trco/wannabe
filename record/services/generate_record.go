package services

import (
	"wannabe/record/actions"
	"wannabe/types"
)

func GenerateRecord(config types.Records, payload types.RecordPayload) ([]byte, error) {
	return actions.GenerateRecord(config, payload)
}
