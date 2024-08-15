package actions

import (
	"encoding/json"
	"fmt"

	"github.com/trco/wannabe/types"
)

func ExtractRecords(requestBody []byte) ([]types.Record, error) {
	var records []types.Record

	err := json.Unmarshal(requestBody, &records)
	if err != nil {
		return nil, fmt.Errorf("ExtractRecords: failed unmarshaling request body: %v", err)
	}

	return records, nil
}
