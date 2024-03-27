package handlers

import (
	"encoding/json"
	"fmt"
	"time"
	"wannabe/config"
	curl "wannabe/curl/services"
	hash "wannabe/hash/services"
	"wannabe/providers"
	"wannabe/record/entities"

	"github.com/gofiber/fiber/v2"
)

func Regenerate(config config.Config, storageProvider providers.StorageProvider) WannabeHandler {
	return func(ctx *fiber.Ctx) error {
		if !config.StorageProvider.Regenerate {
			return internalError(ctx, fmt.Errorf("regenerate set to false in config"))
		}

		regenCount := 0
		regenHashes := []string{}
		failedCount := 0
		failedHashes := []string{}

		hashes, err := storageProvider.GetHashes()
		if err != nil {
			return internalError(ctx, err)
		}

		// REVIEW mem issue in case of too many records ?
		records, err := storageProvider.ReadRecords(hashes)
		if err != nil {
			return internalError(ctx, err)
		}

		for _, recordBytes := range records {
			var record entities.Record

			err := json.Unmarshal(recordBytes, &record)
			oldHash := record.Request.Hash
			if err != nil {
				failedCount++
				failedHashes = append(failedHashes, oldHash)
				continue
			}

			bodyBytes, err := json.Marshal(record.Request.Body)
			if err != nil {
				failedCount++
				failedHashes = append(failedHashes, oldHash)
			}

			curl, err := curl.GenerateCurl(
				record.Request.HttpMethod,
				record.Request.Path,
				record.Request.Query,
				record.Request.Headers,
				bodyBytes,
				config,
			)
			if err != nil {
				failedCount++
				failedHashes = append(failedHashes, oldHash)
			}

			hash, err := hash.GenerateHash(curl)
			if err != nil {
				failedCount++
				failedHashes = append(failedHashes, oldHash)
			}

			isDuplicateHash := checkDuplicates(hashes, hash)
			isDuplicateRegenHash := checkDuplicates(regenHashes, hash)
			if isDuplicateHash || isDuplicateRegenHash {
				continue
			}

			record.Request.Hash = hash
			record.Request.Curl = curl
			record.Metadata.RegeneratedAt = entities.Timestamp{
				Unix: time.Now().Unix(),
				UTC:  time.Now().UTC(),
			}

			recordBytesRegen, err := json.Marshal(record)
			if err != nil {
				failedCount++
				failedHashes = append(failedHashes, oldHash)
			}

			err = storageProvider.InsertRecords([]string{hash}, [][]byte{recordBytesRegen})
			if err != nil {
				failedCount++
				failedHashes = append(failedHashes, oldHash)
			}

			regenCount++
			regenHashes = append(regenHashes, hash)
		}

		return ctx.Status(fiber.StatusCreated).JSON(RegenerateResponse{
			Message:           fmt.Sprintf("%v records succeeded in regenerating, %v records failed in regenerating", regenCount, failedCount),
			RegeneratedHashes: regenHashes,
			FailedHashes:      failedHashes,
		})
	}
}
