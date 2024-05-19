package actions

import (
	"encoding/json"
	"fmt"
	"wannabe/types"
)

func ExtractRecords(requestBody [][]byte) ([]types.Record, error) {
	var records []types.Record

	for _, item := range requestBody {
		var record types.Record
		err := json.Unmarshal(item, &record)

		if err != nil {
			return nil, fmt.Errorf("ExtractRecords: failed unmarshaling request body: %v", err)
		}

		records = append(records, record)

	}

	return records, nil
}
