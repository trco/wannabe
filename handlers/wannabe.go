package handlers

import (
	"fmt"
	"wannabe/config"
	curl "wannabe/curl/services"
	hash "wannabe/hash/services"
	"wannabe/providers"
	record "wannabe/record/services"
	response "wannabe/response/services"

	"github.com/gofiber/fiber/v2"
)

func Wannabe(config config.Config, storageProvider providers.StorageProvider) WannabeHandler {
	return func(ctx *fiber.Ctx) error {
		curl, err := curl.GenerateCurl(
			ctx.Method(),
			ctx.Path(),
			ctx.Queries(),
			ctx.GetReqHeaders(),
			ctx.Body(),
			config,
		)
		if err != nil {
			return internalError(ctx, err)
		}

		hash, err := hash.GenerateHash(curl)
		if err != nil {
			return internalError(ctx, err)
		}

		if config.Read.Enabled {
			records, err := storageProvider.ReadRecords([]string{hash})
			if err != nil && config.Read.FailOnError {
				return internalError(ctx, err)
			}

			if records != nil {
				// response from record is set directly to ctx
				err = response.SetResponse(ctx, records[0])
				if err != nil && config.Read.FailOnError {
					return internalError(ctx, err)
				}

				// TODO remove log
				fmt.Println("GetResponse >>> READ and return")

				return nil
			}
		}

		// response is set directly to ctx
		err = response.FetchResponse(ctx, config.Server)
		if err != nil {
			return internalError(ctx, err)
		}

		record, err := record.GenerateRecord(ctx, config.Records, config.Server, curl, hash)
		if err != nil {
			return internalError(ctx, err)
		}

		err = storageProvider.InsertRecords([]string{hash}, [][]byte{record})
		if err != nil {
			return err
		}

		return nil
	}
}
