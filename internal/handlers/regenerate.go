package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/trco/wannabe/internal/config"
	"github.com/trco/wannabe/internal/hash"
	"github.com/trco/wannabe/internal/record"
	"github.com/trco/wannabe/internal/storage"
)

func Regenerate(cfg config.Config, storageProvider storage.StorageProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetRegenerate(cfg, storageProvider, w, r)
		default:
			InternalErrorApi(w, errors.New("invalid method"), http.StatusMethodNotAllowed)
		}
	}
}

func GetRegenerate(cfg config.Config, storageProvider storage.StorageProvider, w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	if host == "" {
		InternalErrorApi(w, errors.New("required query parameter missing: 'host'"), http.StatusBadRequest)
		return
	}

	regeneratedCount := 0
	regeneratedHashes := []string{}
	failedCount := 0
	failedHashes := []string{}

	hashes, err := storageProvider.GetHashes(host)
	if err != nil {
		InternalErrorApi(w, err, http.StatusInternalServerError)
		return
	}

	// REVIEW could we hit mem issue in case of too many records
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

	for _, r := range records {
		oldHash := r.Request.Hash

		request, err := record.GenerateRequest(r.Request)
		if err != nil {
			failedCount++
			failedHashes = append(failedHashes, oldHash)

			continue
		}

		wannabe := cfg.Wannabes[host]

		curl, err := hash.GenerateCurl(request, wannabe)
		if err != nil {
			failedCount++
			failedHashes = append(failedHashes, oldHash)

			continue
		}

		hash, err := hash.Generate(curl)
		if err != nil {
			failedCount++
			failedHashes = append(failedHashes, oldHash)

			continue
		}

		r.Request.Hash = hash
		r.Request.Curl = curl
		r.Metadata.RegeneratedAt = record.Timestamp{
			Unix: time.Now().Unix(),
			UTC:  time.Now().UTC(),
		}

		encodedRegeneratedRecord, err := json.Marshal(r)
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

	ApiResponse(w, RegenerateResponse{
		Message:           fmt.Sprintf("%v records succeeded in regenerating, %v records failed in regenerating", regeneratedCount, failedCount),
		RegeneratedHashes: regeneratedHashes,
		FailedHashes:      failedHashes,
	})
}

type RegenerateResponse struct {
	Message           string   `json:"message"`
	RegeneratedHashes []string `json:"regeneratedHashes"`
	FailedHashes      []string `json:"failedHashes"`
}
