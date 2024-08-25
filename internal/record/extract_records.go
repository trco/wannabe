package record

import (
	"encoding/json"
	"fmt"
)

func ExtractRecords(requestBody []byte) ([]Record, error) {
	var records []Record

	err := json.Unmarshal(requestBody, &records)
	if err != nil {
		return nil, fmt.Errorf("ExtractRecords: failed unmarshaling request body: %v", err)
	}

	return records, nil
}
