package handler

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

func Regenerate(cfg config.Config, storageProvider storage.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getRegenerate(cfg, storageProvider, w, r)
		default:
			internalErrorApi(w, errors.New("invalid method"), http.StatusMethodNotAllowed)
		}
	}
}

func getRegenerate(cfg config.Config, storageProvider storage.Provider, w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	if host == "" {
		internalErrorApi(w, errors.New("required query parameter missing: 'host'"), http.StatusBadRequest)
		return
	}

	regeneratedCount := 0
	regeneratedHashes := []string{}
	failedCount := 0
	failedHashes := []string{}

	hashes, err := storageProvider.GetHashes(host)
	if err != nil {
		internalErrorApi(w, err, http.StatusInternalServerError)
		return
	}

	// REVIEW could we hit mem issue in case of too many records
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

	for _, rec := range records {
		oldHash := rec.Request.Hash

		request, err := record.GenerateRequest(rec.Request)
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

		rec.Request.Hash = hash
		rec.Request.Curl = curl
		rec.Metadata.RegeneratedAt = record.Timestamp{
			Unix: time.Now().Unix(),
			UTC:  time.Now().UTC(),
		}

		encodedRegeneratedRecord, err := json.Marshal(rec)
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

	apiResponse(w, regenerateResponse{
		Message:           fmt.Sprintf("%v records succeeded in regenerating, %v records failed in regenerating", regeneratedCount, failedCount),
		RegeneratedHashes: regeneratedHashes,
		FailedHashes:      failedHashes,
	})
}

type regenerateResponse struct {
	Message           string   `json:"message"`
	RegeneratedHashes []string `json:"regeneratedHashes"`
	FailedHashes      []string `json:"failedHashes"`
}
