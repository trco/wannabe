package actions

import (
	"encoding/json"
	"fmt"
	"wannabe/record/entities"
)

func ExtractRecords(requestBody [][]byte) ([]entities.Record, error) {
	var records []entities.Record

	for _, item := range requestBody {
		var record entities.Record
		err := json.Unmarshal(item, &record)

		if err != nil {
			return nil, fmt.Errorf("ExtractRecords: failed unmarshaling request body: %v", err)
		}

		records = append(records, record)

	}

	return records, nil
}
