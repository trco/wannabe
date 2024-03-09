package services

import (
	"encoding/json"
	"fmt"
	"wannabe/record/entities"
)

// FIXME add action
func PrepareRecords(recordsBytes [][]byte) ([]entities.Record, error) {
	var records []entities.Record

	for _, recordBytes := range recordsBytes {
		var record entities.Record

		err := json.Unmarshal(recordBytes, &record)
		if err != nil {
			// REVIEW
			// return also corrupted records ?
			// return valid records and corrupted records with errors
			return nil, fmt.Errorf("PrepareRecords: failed unmarshaling record: %v", err)
		}

		records = append(records, record)
	}

	return records, nil
}
