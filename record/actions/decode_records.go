package actions

import (
	"encoding/json"
	"fmt"
	"wannabe/record/entities"
)

func DecodeRecords(encodedRecords [][]byte) ([]entities.Record, error) {
	var records []entities.Record

	for _, encodedRecord := range encodedRecords {
		var record entities.Record

		err := json.Unmarshal(encodedRecord, &record)
		if err != nil {
			// REVIEW
			// return also corrupted records ?
			// return valid records and corrupted records with errors
			return nil, fmt.Errorf("DecodeRecords: failed unmarshaling record: %v", err)
		}

		records = append(records, record)
	}

	return records, nil
}
