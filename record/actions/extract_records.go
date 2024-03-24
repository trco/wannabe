package actions

import (
	"encoding/json"
	"fmt"
	"wannabe/record/common"
)

func ExtractRecords(bodyBytes []byte) ([]common.Record, error) {
	var records []common.Record

	err := json.Unmarshal(bodyBytes, &records)
	if err != nil {
		return nil, fmt.Errorf("ExtractRecords: failed unmarshaling request body: %v", err)
	}

	return records, nil
}
