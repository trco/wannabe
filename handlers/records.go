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
		// TODO validate ctx.Body in relation to config, config.RequestMatching

		records, err := record.ExtractRecords(ctx.Body())
		if err != nil {
			return internalError(ctx, err)
		}

		var hashes []string
		var encodedRecords [][]byte

		for _, record := range records {
			body := record.Request.Body
			var bodyBytes []byte

			if body == nil {
				bodyBytes = []byte("")
			} else {
				bodyBytes, err = json.Marshal(body)
				if err != nil {
					return internalError(ctx, fmt.Errorf("PostRecords: failed marshaling record's request body: %v", err))
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
				return internalError(ctx, err)
			}

			hash, err := hash.GenerateHash(curl)
			if err != nil {
				return internalError(ctx, err)
			}

			if checkDuplicates(hashes, hash) {
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
