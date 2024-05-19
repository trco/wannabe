package services

import (
	"wannabe/record/actions"
	"wannabe/types"
)

func GenerateRecordPayload(wannabeSession types.WannabeSession, hash string, curl string) (types.RecordPayload, error) {
	return actions.GenerateRecordPayload(wannabeSession, hash, curl)
}
