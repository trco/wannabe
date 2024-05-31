package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	curl "wannabe/curl/services"
	"wannabe/handlers/utils"
	hash "wannabe/hash/actions"
	"wannabe/providers"
	recordActions "wannabe/record/actions"
	"wannabe/types"
)

func Regenerate(config types.Config, storageProvider providers.StorageProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetRegenerate(config, storageProvider, w, r)
		default:
			utils.InternalErrorApi(w, errors.New("invalid method"), http.StatusMethodNotAllowed)
		}
	}
}

func GetRegenerate(config types.Config, storageProvider providers.StorageProvider, w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	if host == "" {
		utils.InternalErrorApi(w, errors.New("required query parameter missing: 'host'"), http.StatusBadRequest)
		return
	}

	regeneratedCount := 0
	regeneratedHashes := []string{}
	failedCount := 0
	failedHashes := []string{}

	hashes, err := storageProvider.GetHashes(host)
	if err != nil {
		utils.InternalErrorApi(w, err, http.StatusInternalServerError)
		return
	}

	// REVIEW mem issue in case of too many records ?
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

	for _, record := range records {
		oldHash := record.Request.Hash

		request, err := recordActions.GenerateRequest(record.Request)
		if err != nil {
			failedCount++
			failedHashes = append(failedHashes, oldHash)

			continue
		}

		wannabe := config.Wannabes[host]

		curl, err := curl.GenerateCurl(request, wannabe)
		if err != nil {
			failedCount++
			failedHashes = append(failedHashes, oldHash)

			continue
		}

		hash, err := hash.GenerateHash(curl)
		if err != nil {
			failedCount++
			failedHashes = append(failedHashes, oldHash)

			continue
		}

		record.Request.Hash = hash
		record.Request.Curl = curl
		record.Metadata.RegeneratedAt = types.Timestamp{
			Unix: time.Now().Unix(),
			UTC:  time.Now().UTC(),
		}

		encodedRegeneratedRecord, err := json.Marshal(record)
		if err != nil {
			failedCount++
			failedHashes = append(failedHashes, oldHash)

			continue
		}

		err = storageProvider.InsertRecords(host, []string{hash}, [][]byte{encodedRegeneratedRecord}, true)
		if err != nil {
			failedCount++
			failedHashes = append(failedHashes, oldHash)

			continue
		}

		regeneratedCount++
		regeneratedHashes = append(regeneratedHashes, hash)
	}

	utils.ApiResponse(w, types.RegenerateResponse{
		Message:           fmt.Sprintf("%v records succeeded in regenerating, %v records failed in regenerating", regeneratedCount, failedCount),
		RegeneratedHashes: regeneratedHashes,
		FailedHashes:      failedHashes,
	})
}
