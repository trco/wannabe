package handlers

import (
	"testing"
)

func TestProcessRecordValidation(t *testing.T) {
	t.Run("process record validation", func(t *testing.T) {
		count := 0
		var recordProcessingDetails []RecordProcessingDetails

		ProcessRecordValidation(&recordProcessingDetails, "test hash", "test message", &count)

		wantHash := "test hash"
		wantMessage := "test message"
		wantCount := 1

		gotHash := recordProcessingDetails[0].Hash
		gotMessage := recordProcessingDetails[0].Message
		gotCount := 1

		if gotHash != wantHash {
			t.Errorf("got hash %v, want hash %v", gotHash, wantHash)
		}

		if gotMessage != wantMessage {
			t.Errorf("got message %v, want message %v", gotMessage, wantMessage)
		}

		if gotCount != wantCount {
			t.Errorf("got count %v, want count %v", gotCount, wantCount)
		}
	})
}
