package handlers

import (
	"fmt"
	"wannabe/config"
	curlEntities "wannabe/curl/entities"
	curl "wannabe/curl/services"
	hash "wannabe/hash/services"
	"wannabe/providers"
	recordEntities "wannabe/record/entities"
	record "wannabe/record/services"
	response "wannabe/response/services"

	"github.com/gofiber/fiber/v2"
)

func Wannabe(config config.Config, storageProvider providers.StorageProvider) WannabeHandler {
	return func(ctx *fiber.Ctx) error {
		curlPayload := curlEntities.GenerateCurlPayload{
			HttpMethod:     ctx.Method(),
			Path:           ctx.Path(),
			Query:          ctx.Queries(),
			RequestHeaders: ctx.GetReqHeaders(),
			RequestBody:    ctx.Body(),
		}

		curl, err := curl.GenerateCurl(config, curlPayload)
		if err != nil {
			return internalError(ctx, err)
		}

		hash, err := hash.GenerateHash(curl)
		if err != nil {
			return internalError(ctx, err)
		}

		// server, mixed
		if config.Mode != "proxy" {
			records, err := storageProvider.ReadRecords([]string{hash})
			if err != nil && config.FailOnReadError {
				return internalError(ctx, err)
			}

			if records != nil {
				// response from record is set directly to ctx
				err = response.SetResponse(ctx, records[0])
				if err != nil && config.FailOnReadError {
					return internalError(ctx, err)
				}

				// TODO remove log
				fmt.Println("GetResponse >>> READ and return")

				return nil
			}

			if config.Mode == "server" {
				return internalError(ctx, fmt.Errorf("no record found for the request"))
			}
		}

		// response is set directly to ctx
		err = response.FetchResponse(ctx, config.Server)
		if err != nil {
			return internalError(ctx, err)
		}

		recordPayload := recordEntities.GenerateRecordPayload{
			Hash:            hash,
			Curl:            curl,
			HttpMethod:      ctx.Method(),
			Host:            config.Server,
			Path:            ctx.Path(),
			Query:           ctx.Queries(),
			RequestHeaders:  ctx.GetReqHeaders(),
			RequestBody:     ctx.Body(),
			StatusCode:      ctx.Response().StatusCode(),
			ResponseHeaders: ctx.GetRespHeaders(),
			ResponseBody:    ctx.Response().Body(),
			Timestamp: recordEntities.Timestamp{
				Unix: ctx.Context().Time().Unix(),
				UTC:  ctx.Context().Time().UTC(),
			},
		}

		record, err := record.GenerateRecord(config.Records, recordPayload)
		if err != nil {
			return internalError(ctx, err)
		}

		err = storageProvider.InsertRecords([]string{hash}, [][]byte{record})
		if err != nil {
			return internalError(ctx, err)
		}

		return nil
	}
}
