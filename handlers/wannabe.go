package handlers

import (
	"wannabe/config"
	"wannabe/providers"

	"github.com/gofiber/fiber/v2"
)

func Wannabe(config config.Config, storageProvider providers.StorageProvider) WannabeHandlerObsolete {
	return func(ctx *fiber.Ctx) error {
		// fmt.Println("ctx", ctx)

		// host := strings.Join(ctx.GetReqHeaders()["Host"], "")
		// wannabe := config.Wannabes[host]

		// fmt.Println("host", host)
		// // fmt.Println("wannabe", wannabe)

		// curlPayload := curlEntities.GenerateCurlPayload{
		// 	HttpMethod: ctx.Method(),
		// 	Host:       host,
		// 	Path:       ctx.Path(),
		// 	// Query:          ctx.Queries(),
		// 	// FIXME wrong, but just to circumvent error before removing this handler
		// 	Query:          ctx.GetReqHeaders(),
		// 	RequestHeaders: ctx.GetReqHeaders(),
		// 	RequestBody:    ctx.Body(),
		// }

		// curl, err := curl.GenerateCurl(curlPayload, wannabe)
		// if err != nil {
		// 	return internalError(ctx, err)
		// }

		// hash, err := hash.GenerateHash(curl)
		// if err != nil {
		// 	return internalError(ctx, err)
		// }

		// // server, mixed
		// if config.Mode != "proxy" {
		// 	records, err := storageProvider.ReadRecords([]string{hash}, host)
		// 	if err != nil && config.FailOnReadError {
		// 		return internalError(ctx, err)
		// 	}

		// 	if records != nil {
		// 		// response from record is set directly to ctx
		// 		err = response.SetResponse(ctx, records[0])
		// 		if err != nil && config.FailOnReadError {
		// 			return internalError(ctx, err)
		// 		}

		// 		// TODO remove log
		// 		fmt.Println("GetResponse >>> READ and return")

		// 		return nil
		// 	}

		// 	if config.Mode == "server" {
		// 		return internalError(ctx, fmt.Errorf("no record found for the request"))
		// 	}
		// }

		// // response is set directly to ctx
		// err = response.FetchResponse(ctx, wannabe.Protocol, host)
		// if err != nil {
		// 	return internalError(ctx, err)
		// }

		// recordPayload := recordEntities.GenerateRecordPayload{
		// 	Hash:            hash,
		// 	Curl:            curl,
		// 	HttpMethod:      ctx.Method(),
		// 	Host:            host,
		// 	Path:            ctx.Path(),
		// 	Query:           ctx.Queries(),
		// 	RequestHeaders:  ctx.GetReqHeaders(),
		// 	RequestBody:     ctx.Body(),
		// 	StatusCode:      ctx.Response().StatusCode(),
		// 	ResponseHeaders: ctx.GetRespHeaders(),
		// 	ResponseBody:    ctx.Response().Body(),
		// 	Timestamp: recordEntities.Timestamp{
		// 		Unix: ctx.Context().Time().Unix(),
		// 		UTC:  ctx.Context().Time().UTC(),
		// 	},
		// }

		// record, err := record.GenerateRecord(wannabe.Records, recordPayload)
		// if err != nil {
		// 	return internalError(ctx, err)
		// }

		// err = storageProvider.InsertRecords([][]byte{record}, []string{hash}, host)
		// if err != nil {
		// 	return internalError(ctx, err)
		// }

		return nil
	}
}
