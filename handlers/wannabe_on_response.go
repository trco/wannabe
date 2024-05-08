package handlers

import (
	"log"
	"net/http"
	"wannabe/providers"

	"github.com/AdguardTeam/gomitmproxy"
)

func WannabeOnResponse(storageProvider providers.StorageProvider) WannabeOnResponseHandler {
	return func(session *gomitmproxy.Session) *http.Response {
		log.Printf("onResponse: %s", session.Request().URL.String())

		// FIXME block processing based on sessions.Prop
		if _, blocked := session.GetProp("blocked"); blocked {
			return nil
		}
		if _, responseFromRecord := session.GetProp("responseFromRecord"); responseFromRecord {
			return nil
		}

		// response := session.Response()

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
