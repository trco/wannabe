package services

import (
	"wannabe/record/actions"
	"wannabe/types"

	"github.com/AdguardTeam/gomitmproxy"
)

func GenerateRecordPayload(session *gomitmproxy.Session, hash string, curl string) (types.RecordPayload, error) {
	return actions.GenerateRecordPayload(session, hash, curl)
}
