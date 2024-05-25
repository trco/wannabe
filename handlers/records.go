package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
	curl "wannabe/curl/services"
	hash "wannabe/hash/services"
	"wannabe/providers"
	recordServices "wannabe/record/services"
	"wannabe/types"
)

func Records(config types.Config, storageProvider providers.StorageProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetRecords(storageProvider, w, r)
		case http.MethodDelete:
			DeleteRecords(storageProvider, w, r)
		case http.MethodPost:
			PostRecords(config, storageProvider, w, r)
		default:
			internalErrorApi(w, errors.New("invalid method"), http.StatusMethodNotAllowed)
		}
	}
}

func GetRecords(storageProvider providers.StorageProvider, w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	if host == "" {
		internalErrorApi(w, errors.New("required query parameter missing: 'host'"), http.StatusBadRequest)
		return
	}

	hash := r.PathValue("hash")
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

	records, err := recordServices.DecodeRecords(encodedRecords)
	if err != nil {
		internalErrorApi(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

func DeleteRecords(storageProvider providers.StorageProvider, w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	if host == "" {
		internalErrorApi(w, errors.New("required query parameter missing: 'host'"), http.StatusBadRequest)
		return
	}

	hash := r.PathValue("hash")
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(types.DeleteRecordsResponse{
		Message: fmt.Sprintf("%v records successfully deleted.", len(hashes)),
		Hashes:  hashes,
	})
}

func PostRecords(config types.Config, storageProvider providers.StorageProvider, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		internalErrorApi(w, err, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	records, err := recordServices.ExtractRecords(body)
	if err != nil {
		internalErrorApi(w, err, http.StatusInternalServerError)
		return
	}

	validationErrors, err := recordServices.ValidateRecords(records)
	if err != nil {
		internalErrorApi(w, err, http.StatusInternalServerError)
		return
	}

	insertedCount := 0
	notInsertedCount := 0
	var recordProcessingDetails []types.RecordProcessingDetails

	for i, record := range records {
		if validationErrors[i] != "" {
			processRecordValidation(&recordProcessingDetails, "", validationErrors[i], &notInsertedCount)
			continue
		}

		request, err := recordServices.GenerateRequest(record.Request)
		if err != nil {
			processRecordValidation(&recordProcessingDetails, "", err.Error(), &notInsertedCount)
			continue
		}

		host := record.Request.Host
		wannabe := config.Wannabes[host]

		curl, err := curl.GenerateCurl(request, wannabe)
		if err != nil {
			processRecordValidation(&recordProcessingDetails, "", err.Error(), &notInsertedCount)
			continue
		}

		hash, err := hash.GenerateHash(curl)
		if err != nil {
			processRecordValidation(&recordProcessingDetails, "", err.Error(), &notInsertedCount)
			continue
		}

		record.Request.Curl = curl
		record.Request.Hash = hash
		record.Metadata.GeneratedAt = types.Timestamp{
			Unix: time.Now().Unix(),
			UTC:  time.Now().UTC(),
		}

		encodedRecord, err := json.Marshal(record)
		if err != nil {
			processRecordValidation(&recordProcessingDetails, hash, err.Error(), &notInsertedCount)
			continue
		}

		err = storageProvider.InsertRecords(host, []string{hash}, [][]byte{encodedRecord})
		if err != nil {
			internalErrorApi(w, err, http.StatusInternalServerError)
			return
		}

		processRecordValidation(&recordProcessingDetails, hash, "success", &insertedCount)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(types.PostRecordsResponse{
		InsertedRecordsCount:    insertedCount,
		NotInsertedRecordsCount: notInsertedCount,
		RecordProcessingDetails: recordProcessingDetails,
	})
}
