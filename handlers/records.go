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

		validationErrors, err := record.ValidateRecords(config, records)
		if err != nil {
			return internalError(ctx, err)
		}

		insertedCount := 0
		notInsertedCount := 0
		var recordProcessingDetails []RecordProcessingDetails

		var validHashes []string
		var encodedRecords [][]byte

		for i, record := range records {
			if validationErrors[i] != "" {
				processRecordValidation(&recordProcessingDetails, "", validationErrors[i], &notInsertedCount)

				continue
			}

			body := record.Request.Body
			var requestBody []byte

			if body == nil {
				requestBody = []byte("")
			} else {
				requestBody, err = json.Marshal(body)
				if err != nil {
					processRecordValidation(&recordProcessingDetails, "", err.Error(), &notInsertedCount)

					continue
				}
			}

			curl, err := curl.GenerateCurl(
				record.Request.HttpMethod,
				record.Request.Path,
				record.Request.Query,
				record.Request.Headers,
				requestBody,
				config,
			)
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
			record.Metadata.GeneratedAt = entities.Timestamp{
				Unix: time.Now().Unix(),
				UTC:  time.Now().UTC(),
			}

			encodedRecord, err := json.Marshal(record)
			if err != nil {
				processRecordValidation(&recordProcessingDetails, hash, err.Error(), &notInsertedCount)

				continue
			}

			validHashes = append(validHashes, hash)
			encodedRecords = append(encodedRecords, encodedRecord)

			processRecordValidation(&recordProcessingDetails, hash, "success", &insertedCount)
		}

		err = storageProvider.InsertRecords(validHashes, encodedRecords)
		if err != nil {
			return internalError(ctx, err)
		}

		return ctx.Status(fiber.StatusCreated).JSON(PostRecordsResponse{
			InsertedRecordsCount:    insertedCount,
			NotInsertedRecordsCount: notInsertedCount,
			RecordProcessingDetails: recordProcessingDetails,
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
