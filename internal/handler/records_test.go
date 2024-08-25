package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInternalErrorApi(t *testing.T) {
	t.Run("internal error api", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := &TestError{"test error"}
		status := 500

		internalErrorApi(w, err, status)
		responseBody := w.Body.String()

		wantStatusCode := http.StatusInternalServerError
		wantContentType := "application/json"
		wantResponseBody := "{\"error\":\"test error\"}\n"

		if w.Code != wantStatusCode {
			t.Errorf("got status code %v, want status code %v", w.Code, wantStatusCode)
		}

		if w.Header().Get("Content-Type") != wantContentType {
			t.Errorf("got content type %v, want content type %v", w.Header().Get("Content-Type"), wantContentType)
		}

		if responseBody != wantResponseBody {
			t.Errorf("got response body %v, want response body %v", responseBody, wantResponseBody)
		}
	})
}

func TestApiResponse(t *testing.T) {
	t.Run("api response", func(t *testing.T) {
		w := httptest.NewRecorder()
		response := map[string]string{"key": "value"}

		apiResponse(w, response)
		responseBody := w.Body.String()

		wantStatusCode := http.StatusOK
		wantContentType := "application/json"
		wantResponseBody := "{\"key\":\"value\"}\n"

		if w.Code != wantStatusCode {
			t.Errorf("got status code %v, want status code %v", w.Code, wantStatusCode)
		}

		if w.Header().Get("Content-Type") != wantContentType {
			t.Errorf("got content type %v, want content type %v", w.Header().Get("Content-Type"), wantContentType)
		}

		if responseBody != wantResponseBody {
			t.Errorf("got response body %v, want response body %v", responseBody, wantResponseBody)
		}
	})
}

func TestProcessRecordValidation(t *testing.T) {
	t.Run("process record validation", func(t *testing.T) {
		count := 0
		var recordProcessingDetails []recordProcessingDetails

		processRecordValidation(&recordProcessingDetails, "test hash", "test message", &count)

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
