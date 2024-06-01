package utils

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"wannabe/types"
)

func TestInternalErrorOnRequest(t *testing.T) {
	session := &MockSession{
		props: make(map[string]interface{}),
	}
	request, _ := http.NewRequest("GET", "test.com", nil)

	_, response := InternalErrorOnRequest(session, request, &TestError{"test error"})

	if http.StatusInternalServerError != response.StatusCode {
		t.Errorf("expected response status code: %v, actual response status code: %v", http.StatusInternalServerError, response.StatusCode)
	}

	contentType := response.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type header: %v, actual Content-Type header %v", "application/json", contentType)
	}

	body, _ := io.ReadAll(response.Body)
	responseBody := string(body)

	expectedResponseBody := "{\"error\":\"test error\"}"

	if expectedResponseBody != responseBody {
		t.Errorf("expected response body: %v, actual response body: %v", expectedResponseBody, responseBody)
	}
}

func TestInternalErrorOnResponse(t *testing.T) {
	request, _ := http.NewRequest("GET", "http://test.com", nil)
	response := InternalErrorOnResponse(request, &TestError{"test error"})

	if http.StatusInternalServerError != response.StatusCode {
		t.Errorf("expected response status code: %v, actual response status code: %v", http.StatusInternalServerError, response.StatusCode)
	}

	contentType := response.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type header: %v, actual Content-Type header %v", "application/json", contentType)
	}

	body, _ := io.ReadAll(response.Body)
	responseBody := string(body)

	expectedResponseBody := "{\"error\":\"test error\"}"

	if expectedResponseBody != responseBody {
		t.Errorf("expected response body: %v, actual response body: %v", expectedResponseBody, responseBody)
	}
}

func TestPrepareResponseBody(t *testing.T) {
	reader := PrepareResponseBody(&TestError{"test error"})

	responseBody, _ := io.ReadAll(reader)

	expected := "{\"error\":\"test error\"}"

	if !bytes.Equal([]byte(expected), responseBody) {
		t.Errorf("expected response body: %v, actual response body: %v", expected, string(responseBody))
	}
}

func TestShouldSkipResponseProcessing(t *testing.T) {
	sessionBlocked := &MockSession{
		props: make(map[string]interface{}),
	}

	skip := ShouldSkipResponseProcessing(sessionBlocked)

	if skip != false {
		t.Errorf("expected skip: false, actual skip: %v, when 'blocked' and 'responseSetFromRecord' are not set on session", skip)
	}

	sessionBlocked.SetProp("blocked", true)

	skip = ShouldSkipResponseProcessing(sessionBlocked)

	if skip != true {
		t.Errorf("expected skip: true, actual skip: %v, when 'blocked' is set to 'true' on session", skip)
	}

	sessionResponseSetFromRecord := &MockSession{
		props: make(map[string]interface{}),
	}

	sessionResponseSetFromRecord.SetProp("responseSetFromRecord", true)

	skip = ShouldSkipResponseProcessing(sessionResponseSetFromRecord)

	if skip != true {
		t.Errorf("expected skip: true, actual skip: %v, when 'responseSetFromRecord' is set to 'true' on session", skip)
	}
}

type TestCaseGetHashAndCurlFromSession struct {
	Session *MockSession
	Hash    string
	Curl    string
}

func TestGetHashAndCurlFromSession(t *testing.T) {
	testCases := map[string]TestCaseGetHashAndCurlFromSession{
		"noHash": {
			Session: &MockSession{
				props: map[string]interface{}{
					"noHash": "test hash",
				},
			},
			Hash: "",
		},
		"hashNotString": {
			Session: &MockSession{
				props: map[string]interface{}{
					"hash": true,
				},
			},
			Hash: "",
		},
		"noCurl": {
			Session: &MockSession{
				props: map[string]interface{}{
					"hash":   "test hash",
					"noCurl": "test curl",
				},
			},
			Hash: "",
			Curl: "",
		},
		"curlNotString": {
			Session: &MockSession{
				props: map[string]interface{}{
					"hash": "test hash",
					"curl": true,
				},
			},
			Hash: "",
			Curl: "",
		},
		"noHashAndCurl": {
			Session: &MockSession{
				props: map[string]interface{}{
					"noHash": "test hash",
					"noCurl": "test curl",
				},
			},
			Hash: "",
			Curl: "",
		},
		"hashAndCurl": {
			Session: &MockSession{
				props: map[string]interface{}{
					"hash": "test hash",
					"curl": "test curl",
				},
			},
			Hash: "test hash",
			Curl: "test curl",
		},
	}

	for testName, tc := range testCases {
		hash, curl, _ := GetHashAndCurlFromSession(tc.Session)

		if hash != tc.Hash || curl != tc.Curl {
			t.Errorf("failed test case: %v, expected hash: %v, actual hash: %v, expected curl: %v, actual curl: %v", testName, tc.Hash, hash, tc.Curl, curl)
		}
	}
}

func TestInternalErrorApi(t *testing.T) {
	w := httptest.NewRecorder()
	err := &TestError{"test error"}
	status := 500

	InternalErrorApi(w, err, status)
	responseBody := w.Body.String()
	responseStatusCode := w.Code

	expectedResponseBody := "{\"error\":\"test error\"}\n"

	if http.StatusInternalServerError != responseStatusCode {
		t.Errorf("expected response status code: %v, actual response status code: %v", http.StatusInternalServerError, responseStatusCode)
	}

	if expectedResponseBody != responseBody {
		t.Errorf("expected response body: %v, actual response body: %v", expectedResponseBody, responseBody)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type header: %v, actual Content-Type header %v", "application/json", contentType)
	}
}

func TestApiResponse(t *testing.T) {
	w := httptest.NewRecorder()
	response := map[string]string{"key": "value"}

	ApiResponse(w, response)
	responseBody := w.Body.String()

	expectedResponseBody := "{\"key\":\"value\"}\n"

	if expectedResponseBody != responseBody {
		t.Errorf("expected response body: %v, actual response body: %v", expectedResponseBody, responseBody)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type header: %v, actual Content-Type header %v", "application/json", contentType)
	}
}

func TestCheckDuplicates(t *testing.T) {
	// duplicates exist
	slice := []string{"a", "b", "c"}

	duplicatesExist := CheckDuplicates(slice, "a")

	if duplicatesExist != true {
		t.Errorf("no duplicates detected although they are present in slice")
	}

	// no duplicates
	duplicatesExist = CheckDuplicates(slice, "d")

	if duplicatesExist == true {
		t.Errorf("duplicates detected although they are not present")
	}
}

func TestProcessRecordValidation(t *testing.T) {
	count := 0
	var recordProcessingDetails []types.RecordProcessingDetails

	ProcessRecordValidation(&recordProcessingDetails, "test hash", "test message", &count)

	if recordProcessingDetails[0].Hash != "test hash" || recordProcessingDetails[0].Message != "test message" || count != 1 {
		t.Errorf("record processing details not valid, expected: hash: 'test hash', message: 'test message', count: 1, actual: hash: '%v', message: '%v', count: %v", recordProcessingDetails[0].Hash, recordProcessingDetails[0].Message, count)
	}
}

type MockSession struct {
	props map[string]interface{}
}

func (m *MockSession) SetProp(key string, value interface{}) {
	m.props[key] = value
}

func (m *MockSession) GetProp(key string) (interface{}, bool) {
	v, ok := m.props[key]
	return v, ok
}

type TestError struct {
	message string
}

func (e *TestError) Error() string {
	return e.message
}
