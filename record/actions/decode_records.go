package actions

import (
	"encoding/json"
	"fmt"

	"github.com/trco/wannabe/types"
)

func DecodeRecords(encodedRecords [][]byte) ([]types.Record, error) {
	var records []types.Record

	for _, encodedRecord := range encodedRecords {
		var record types.Record

		err := json.Unmarshal(encodedRecord, &record)
		if err != nil {
			return nil, fmt.Errorf("DecodeRecords: failed unmarshaling record: %v", err)
		}

		records = append(records, record)
	}

	return records, nil
}
