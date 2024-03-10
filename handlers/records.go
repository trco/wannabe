package handlers

import (
	"encoding/json"
	"fmt"
	"wannabe/config"
	curl "wannabe/curl/services"
	hash "wannabe/hash/services"
	"wannabe/providers"
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

		recordsBytes, err := storageProvider.ReadRecords(hashes)
		if err != nil {
			return internalError(ctx, err)
		}

		records, err := record.DecodeRecords(recordsBytes)
		if err != nil {
			return internalError(ctx, err)
		}

		return ctx.Status(fiber.StatusOK).JSON(records)
	}
}

func PostRecords(config config.Config, storageProvider providers.StorageProvider) WannabeHandler {
	return func(ctx *fiber.Ctx) error {
		// TODO validate ctx.Body in relation to config, config.RequestMatching

		records, err := record.ExtractRecords(ctx.Body())
		if err != nil {
			return internalError(ctx, err)
		}

		var hashes []string
		var encodedRecords [][]byte

		for _, record := range records {
			body := record.Request.Body
			bodyBytes, err := json.Marshal(body)
			if err != nil {
				return internalError(ctx, fmt.Errorf("PostRecords: failed marshaling record's request body: %v", err))
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
				return internalError(ctx, err)
			}

			hash, err := hash.GenerateHash(curl)
			if err != nil {
				return internalError(ctx, err)
			}

			record.Request.Curl = curl
			record.Request.Hash = hash

			encodedRecord, err := json.Marshal(record)
			if err != nil {
				return internalError(ctx, fmt.Errorf("PostRecords: failed marshaling record: %v", err))
			}

			hashes = append(hashes, hash)
			encodedRecords = append(encodedRecords, encodedRecord)
		}

		err = storageProvider.InsertRecords(hashes, encodedRecords)
		if err != nil {
			return internalError(ctx, err)
		}

		return ctx.Status(fiber.StatusCreated).JSON(PostRecordsResponse{
			Message: "Records successfully created.",
			Hashes:  hashes,
		})
	}
}

func DeleteRecords(ctx *fiber.Ctx) error {
	// REVIEW ? bulk delete files using goroutines and channels
	return nil
}
