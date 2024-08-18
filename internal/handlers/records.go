package handlers

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

func Records(cfg config.Config, storageProvider storage.StorageProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetRecords(storageProvider, w, r)
		case http.MethodDelete:
			DeleteRecords(storageProvider, w, r)
		case http.MethodPost:
			PostRecords(cfg, storageProvider, w, r)
		default:
			InternalErrorApi(w, errors.New("invalid method"), http.StatusMethodNotAllowed)
		}
	}
}

func GetRecords(storageProvider storage.StorageProvider, w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	hash := r.PathValue("hash")

	if host == "" && hash == "" {
		hostsAndHashes, err := storageProvider.GetHostsAndHashes()
		if err != nil {
			InternalErrorApi(w, err, http.StatusInternalServerError)
			return
		}

		ApiResponse(w, hostsAndHashes)
		return
	}

	if host == "" {
		InternalErrorApi(w, errors.New("required query parameter missing: 'host'"), http.StatusBadRequest)
		return
	}

	hashes := []string{hash}
	if hash == "" {
		var err error
		hashes, err = storageProvider.GetHashes(host)
		if err != nil {
			InternalErrorApi(w, err, http.StatusInternalServerError)
			return
		}
	}

	encodedRecords, err := storageProvider.ReadRecords(host, hashes)
	if err != nil {
		InternalErrorApi(w, err, http.StatusInternalServerError)
		return
	}

	records, err := record.DecodeRecords(encodedRecords)
	if err != nil {
		InternalErrorApi(w, err, http.StatusInternalServerError)
		return
	}

	ApiResponse(w, records)
}

func DeleteRecords(storageProvider storage.StorageProvider, w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	hash := r.PathValue("hash")

	if host == "" {
		InternalErrorApi(w, errors.New("required query parameter missing: 'host'"), http.StatusBadRequest)
		return
	}

	hashes := []string{hash}
	if hash == "" {
		var err error
		hashes, err = storageProvider.GetHashes(host)
		if err != nil {
			InternalErrorApi(w, err, http.StatusInternalServerError)
			return
		}
	}

	err := storageProvider.DeleteRecords(host, hashes)
	if err != nil {
		InternalErrorApi(w, err, http.StatusInternalServerError)
		return
	}

	ApiResponse(w, DeleteRecordsResponse{
		Message: fmt.Sprintf("%v records successfully deleted.", len(hashes)),
		Hashes:  hashes,
	})
}

func PostRecords(cfg config.Config, storageProvider storage.StorageProvider, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		InternalErrorApi(w, err, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	records, err := record.ExtractRecords(body)
	if err != nil {
		InternalErrorApi(w, err, http.StatusInternalServerError)
		return
	}

	validationErrors, err := record.ValidateRecords(records)
	if err != nil {
		InternalErrorApi(w, err, http.StatusInternalServerError)
		return
	}

	insertedCount := 0
	notInsertedCount := 0
	var recordProcessingDetails []RecordProcessingDetails

	for index, rec := range records {
		if validationErrors[index] != "" {
			ProcessRecordValidation(&recordProcessingDetails, "", validationErrors[index], &notInsertedCount)
			continue
		}

		request, err := record.GenerateRequest(rec.Request)
		if err != nil {
			ProcessRecordValidation(&recordProcessingDetails, "", err.Error(), &notInsertedCount)
			continue
		}

		host := rec.Request.Host
		wannabe := cfg.Wannabes[host]

		curl, err := hash.GenerateCurl(request, wannabe)
		if err != nil {
			ProcessRecordValidation(&recordProcessingDetails, "", err.Error(), &notInsertedCount)
			continue
		}

		hash, err := hash.Generate(curl)
		if err != nil {
			ProcessRecordValidation(&recordProcessingDetails, "", err.Error(), &notInsertedCount)
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
			ProcessRecordValidation(&recordProcessingDetails, hash, err.Error(), &notInsertedCount)
			continue
		}

		err = storageProvider.InsertRecords(host, []string{hash}, [][]byte{encodedRecord}, false)
		if err != nil {
			InternalErrorApi(w, err, http.StatusInternalServerError)
			return
		}

		ProcessRecordValidation(&recordProcessingDetails, hash, "success", &insertedCount)
	}

	ApiResponse(w, PostRecordsResponse{
		InsertedRecordsCount:    insertedCount,
		NotInsertedRecordsCount: notInsertedCount,
		RecordProcessingDetails: recordProcessingDetails,
	})
}

func InternalErrorApi(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(InternalError{Error: err.Error()})
}

func ApiResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type DeleteRecordsResponse struct {
	Message string   `json:"message"`
	Hashes  []string `json:"hashes"`
}

type PostRecordsResponse struct {
	InsertedRecordsCount    int                       `json:"insertedRecordsCount"`
	NotInsertedRecordsCount int                       `json:"notInsertedRecordsCount"`
	RecordProcessingDetails []RecordProcessingDetails `json:"recordProcessingDetails"`
}

type RecordProcessingDetails struct {
	Hash    string `json:"hash"`
	Message string `json:"message"`
}

func ProcessRecordValidation(recordProcessingDetails *[]RecordProcessingDetails, hash string, message string, valueToIncrement *int) {
	*recordProcessingDetails = append(*recordProcessingDetails, RecordProcessingDetails{
		Hash:    hash,
		Message: message,
	})

	*valueToIncrement++
}
