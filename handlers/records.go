package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
	curl "wannabe/curl/services"
	"wannabe/handlers/utils"
	hash "wannabe/hash/actions"
	"wannabe/providers"
	recordActions "wannabe/record/actions"
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
			utils.InternalErrorApi(w, errors.New("invalid method"), http.StatusMethodNotAllowed)
		}
	}
}

func GetRecords(storageProvider providers.StorageProvider, w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	hash := r.PathValue("hash")

	if host == "" && hash == "" {
		hostsAndHashes, err := storageProvider.GetHostsAndHashes()
		if err != nil {
			utils.InternalErrorApi(w, err, http.StatusInternalServerError)
			return
		}

		utils.ApiResponse(w, hostsAndHashes)
		return
	}

	if host == "" {
		utils.InternalErrorApi(w, errors.New("required query parameter missing: 'host'"), http.StatusBadRequest)
		return
	}

	hashes := []string{hash}
	if hash == "" {
		var err error
		hashes, err = storageProvider.GetHashes(host)
		if err != nil {
			utils.InternalErrorApi(w, err, http.StatusInternalServerError)
			return
		}
	}

	encodedRecords, err := storageProvider.ReadRecords(host, hashes)
	if err != nil {
		utils.InternalErrorApi(w, err, http.StatusInternalServerError)
		return
	}

	records, err := recordActions.DecodeRecords(encodedRecords)
	if err != nil {
		utils.InternalErrorApi(w, err, http.StatusInternalServerError)
		return
	}

	utils.ApiResponse(w, records)
}

func DeleteRecords(storageProvider providers.StorageProvider, w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	hash := r.PathValue("hash")

	if host == "" {
		utils.InternalErrorApi(w, errors.New("required query parameter missing: 'host'"), http.StatusBadRequest)
		return
	}

	hashes := []string{hash}
	if hash == "" {
		var err error
		hashes, err = storageProvider.GetHashes(host)
		if err != nil {
			utils.InternalErrorApi(w, err, http.StatusInternalServerError)
			return
		}
	}

	err := storageProvider.DeleteRecords(host, hashes)
	if err != nil {
		utils.InternalErrorApi(w, err, http.StatusInternalServerError)
		return
	}

	utils.ApiResponse(w, types.DeleteRecordsResponse{
		Message: fmt.Sprintf("%v records successfully deleted.", len(hashes)),
		Hashes:  hashes,
	})
}

func PostRecords(config types.Config, storageProvider providers.StorageProvider, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.InternalErrorApi(w, err, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	records, err := recordActions.ExtractRecords(body)
	if err != nil {
		utils.InternalErrorApi(w, err, http.StatusInternalServerError)
		return
	}

	validationErrors, err := recordActions.ValidateRecords(records)
	if err != nil {
		utils.InternalErrorApi(w, err, http.StatusInternalServerError)
		return
	}

	insertedCount := 0
	notInsertedCount := 0
	var recordProcessingDetails []types.RecordProcessingDetails

	for i, record := range records {
		if validationErrors[i] != "" {
			utils.ProcessRecordValidation(&recordProcessingDetails, "", validationErrors[i], &notInsertedCount)
			continue
		}

		request, err := recordActions.GenerateRequest(record.Request)
		if err != nil {
			utils.ProcessRecordValidation(&recordProcessingDetails, "", err.Error(), &notInsertedCount)
			continue
		}

		host := record.Request.Host
		wannabe := config.Wannabes[host]

		curl, err := curl.GenerateCurl(request, wannabe)
		if err != nil {
			utils.ProcessRecordValidation(&recordProcessingDetails, "", err.Error(), &notInsertedCount)
			continue
		}

		hash, err := hash.GenerateHash(curl)
		if err != nil {
			utils.ProcessRecordValidation(&recordProcessingDetails, "", err.Error(), &notInsertedCount)
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
			utils.ProcessRecordValidation(&recordProcessingDetails, hash, err.Error(), &notInsertedCount)
			continue
		}

		err = storageProvider.InsertRecords(host, []string{hash}, [][]byte{encodedRecord}, false)
		if err != nil {
			utils.InternalErrorApi(w, err, http.StatusInternalServerError)
			return
		}

		utils.ProcessRecordValidation(&recordProcessingDetails, hash, "success", &insertedCount)
	}

	utils.ApiResponse(w, types.PostRecordsResponse{
		InsertedRecordsCount:    insertedCount,
		NotInsertedRecordsCount: notInsertedCount,
		RecordProcessingDetails: recordProcessingDetails,
	})
}
