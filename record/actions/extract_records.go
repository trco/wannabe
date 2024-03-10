package actions

import (
	"encoding/json"
	"fmt"
	"wannabe/record/entities"
)

func ExtractRecords(body []byte) ([]entities.Record, error) {
	var records []entities.Record

	err := json.Unmarshal(body, &records)
	if err != nil {
		return nil, fmt.Errorf("ExtractRecords: failed unmarshaling request body: %v", err)
	}

	return records, nil
}
