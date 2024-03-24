package handlers

import (
	"encoding/json"
	"fmt"
	"time"
	"wannabe/config"
	curl "wannabe/curl/services"
	hash "wannabe/hash/services"
	"wannabe/providers"
	"wannabe/record/common"
	record "wannabe/record/services"

	"github.com/gofiber/fiber/v2"
)

func GetRecords(storageProvider providers.StorageProvider) WannabeHandler {
	return func(ctx *fiber.Ctx) error {
		hash := ctx.Params("hash")
		hashes := []string{hash}

		if hash == "" {
			var err error
			hashes, err = storageProvider.GetHashes()
			if err != nil {
				return internalError(ctx, err)
			}
		}

		encodedRecords, err := storageProvider.ReadRecords(hashes)
		if err != nil {
			return internalError(ctx, err)
		}

		records, err := record.DecodeRecords(encodedRecords)
		if err != nil {
			return internalError(ctx, err)
		}

		return ctx.Status(fiber.StatusOK).JSON(records)
	}
}

func PostRecords(config config.Config, storageProvider providers.StorageProvider) WannabeHandler {
	return func(ctx *fiber.Ctx) error {
		records, err := record.ExtractRecords(ctx.Body())
		if err != nil {
			return internalError(ctx, err)
		}

		validations, err := record.ValidateRecords(config, records)
		if err != nil {
			return internalError(ctx, err)
		}

		insertedRecordsCount := 0
		notInsertedRecordsCount := 0
		var processingDetails []ProcessingDetails

		var validHashes []string
		var encodedRecords [][]byte

		for i, record := range records {
			if !validations[i].Valid {
				handleRecordProcessing(&processingDetails, "", validations[i].Error, &notInsertedRecordsCount)

				continue
			}

			body := record.Request.Body
			var bodyBytes []byte

			if body == nil {
				bodyBytes = []byte("")
			} else {
				bodyBytes, err = json.Marshal(body)
				if err != nil {
					handleRecordProcessing(&processingDetails, "", err.Error(), &notInsertedRecordsCount)

					continue
				}
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
				handleRecordProcessing(&processingDetails, "", err.Error(), &notInsertedRecordsCount)

				continue
			}

			hash, err := hash.GenerateHash(curl)
			if err != nil {
				handleRecordProcessing(&processingDetails, "", err.Error(), &notInsertedRecordsCount)

				continue
			}

			record.Request.Curl = curl
			record.Request.Hash = hash
			record.Metadata.GeneratedAt = common.Timestamp{
				Unix: time.Now().Unix(),
				UTC:  time.Now().UTC(),
			}

			encodedRecord, err := json.Marshal(record)
			if err != nil {
				handleRecordProcessing(&processingDetails, hash, err.Error(), &notInsertedRecordsCount)

				continue
			}

			validHashes = append(validHashes, hash)
			encodedRecords = append(encodedRecords, encodedRecord)

			handleRecordProcessing(&processingDetails, hash, "success", &insertedRecordsCount)
		}

		err = storageProvider.InsertRecords(validHashes, encodedRecords)
		if err != nil {
			return internalError(ctx, err)
		}

		return ctx.Status(fiber.StatusCreated).JSON(PostRecordsResponse{
			InsertedRecordsCount:    insertedRecordsCount,
			NotInsertedRecordsCount: notInsertedRecordsCount,
			ProcessingDetails:       processingDetails,
		})
	}
}

func DeleteRecords(storageProvider providers.StorageProvider) WannabeHandler {
	return func(ctx *fiber.Ctx) error {
		hash := ctx.Params("hash")
		hashes := []string{hash}

		if hash == "" {
			var err error
			hashes, err = storageProvider.GetHashes()
			if err != nil {
				return internalError(ctx, err)
			}
		}

		err := storageProvider.DeleteRecords(hashes)
		if err != nil {
			return internalError(ctx, err)
		}

		return ctx.Status(fiber.StatusCreated).JSON(DeleteRecordsResponse{
			Message: fmt.Sprintf("%v records successfully deleted.", len(hashes)),
			Hashes:  hashes,
		})
	}
}
