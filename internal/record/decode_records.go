package record

import (
	"encoding/json"
	"fmt"
)

func DecodeRecords(encodedRecords [][]byte) ([]Record, error) {
	var records []Record

	for _, encodedRecord := range encodedRecords {
		var record Record

		err := json.Unmarshal(encodedRecord, &record)
		if err != nil {
			return nil, fmt.Errorf("DecodeRecords: failed unmarshaling record: %v", err)
		}

		records = append(records, record)
	}

	return records, nil
}
