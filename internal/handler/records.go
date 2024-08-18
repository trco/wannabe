package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/trco/wannabe/internal/config"
	"github.com/trco/wannabe/internal/hash"
	"github.com/trco/wannabe/internal/record"
	"github.com/trco/wannabe/internal/storage"
)

func Records(cfg config.Config, storageProvider storage.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getRecords(storageProvider, w, r)
		case http.MethodDelete:
			deleteRecords(storageProvider, w, r)
		case http.MethodPost:
			postRecords(cfg, storageProvider, w, r)
		default:
			internalErrorApi(w, errors.New("invalid method"), http.StatusMethodNotAllowed)
		}
	}
}

func getRecords(storageProvider storage.Provider, w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	hash := r.PathValue("hash")

	if host == "" && hash == "" {
		hostsAndHashes, err := storageProvider.GetHostsAndHashes()
		if err != nil {
			internalErrorApi(w, err, http.StatusInternalServerError)
			return
		}

		apiResponse(w, hostsAndHashes)
		return
	}

	if host == "" {
		internalErrorApi(w, errors.New("required query parameter missing: 'host'"), http.StatusBadRequest)
		return
	}

	hashes := []string{hash}
	if hash == "" {
		var err error
		hashes, err = storageProvider.GetHashes(host)
		if err != nil {
			internalErrorApi(w, err, http.StatusInternalServerError)
			return
		}
	}

	encodedRecords, err := storageProvider.ReadRecords(host, hashes)
	if err != nil {
		internalErrorApi(w, err, http.StatusInternalServerError)
		return
	}

	records, err := record.DecodeRecords(encodedRecords)
	if err != nil {
		internalErrorApi(w, err, http.StatusInternalServerError)
		return
	}

	apiResponse(w, records)
}

func deleteRecords(storageProvider storage.Provider, w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	hash := r.PathValue("hash")

	if host == "" {
		internalErrorApi(w, errors.New("required query parameter missing: 'host'"), http.StatusBadRequest)
		return
	}

	hashes := []string{hash}
	if hash == "" {
		var err error
		hashes, err = storageProvider.GetHashes(host)
		if err != nil {
			internalErrorApi(w, err, http.StatusInternalServerError)
			return
		}
	}

	err := storageProvider.DeleteRecords(host, hashes)
	if err != nil {
		internalErrorApi(w, err, http.StatusInternalServerError)
		return
	}

	apiResponse(w, deleteRecordsResponse{
		Message: fmt.Sprintf("%v records successfully deleted.", len(hashes)),
		Hashes:  hashes,
	})
}

func postRecords(cfg config.Config, storageProvider storage.Provider, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		internalErrorApi(w, err, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	records, err := record.ExtractRecords(body)
	if err != nil {
		internalErrorApi(w, err, http.StatusInternalServerError)
		return
	}

	validationErrors, err := record.ValidateRecords(records)
	if err != nil {
		internalErrorApi(w, err, http.StatusInternalServerError)
		return
	}

	insertedCount := 0
	notInsertedCount := 0
	var recordProcessingDetails []recordProcessingDetails

	for index, rec := range records {
		if validationErrors[index] != "" {
			processRecordValidation(&recordProcessingDetails, "", validationErrors[index], &notInsertedCount)
			continue
		}

		request, err := record.GenerateRequest(rec.Request)
		if err != nil {
			processRecordValidation(&recordProcessingDetails, "", err.Error(), &notInsertedCount)
			continue
		}

		host := rec.Request.Host
		wannabe := cfg.Wannabes[host]

		curl, err := hash.GenerateCurl(request, wannabe)
		if err != nil {
			processRecordValidation(&recordProcessingDetails, "", err.Error(), &notInsertedCount)
			continue
		}

		hash, err := hash.Generate(curl)
		if err != nil {
			processRecordValidation(&recordProcessingDetails, "", err.Error(), &notInsertedCount)
			continue
		}

		rec.Request.Curl = curl
		rec.Request.Hash = hash
		rec.Metadata.GeneratedAt = record.Timestamp{
			Unix: time.Now().Unix(),
			UTC:  time.Now().UTC(),
		}

		encodedRecord, err := json.Marshal(rec)
		if err != nil {
			processRecordValidation(&recordProcessingDetails, hash, err.Error(), &notInsertedCount)
			continue
		}

		err = storageProvider.InsertRecords(host, []string{hash}, [][]byte{encodedRecord}, false)
		if err != nil {
			internalErrorApi(w, err, http.StatusInternalServerError)
			return
		}

		processRecordValidation(&recordProcessingDetails, hash, "success", &insertedCount)
	}

	apiResponse(w, postRecordsResponse{
		InsertedRecordsCount:    insertedCount,
		NotInsertedRecordsCount: notInsertedCount,
		RecordProcessingDetails: recordProcessingDetails,
	})
}

func apiResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func internalErrorApi(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(InternalError{Error: err.Error()})
}

type deleteRecordsResponse struct {
	Message string   `json:"message"`
	Hashes  []string `json:"hashes"`
}

type recordProcessingDetails struct {
	Hash    string `json:"hash"`
	Message string `json:"message"`
}

func processRecordValidation(recProcessingDetails *[]recordProcessingDetails, hash string, message string, valueToIncrement *int) {
	*recProcessingDetails = append(*recProcessingDetails, recordProcessingDetails{
		Hash:    hash,
		Message: message,
	})

	*valueToIncrement++
}

type postRecordsResponse struct {
	InsertedRecordsCount    int                       `json:"insertedRecordsCount"`
	NotInsertedRecordsCount int                       `json:"notInsertedRecordsCount"`
	RecordProcessingDetails []recordProcessingDetails `json:"recordProcessingDetails"`
}
