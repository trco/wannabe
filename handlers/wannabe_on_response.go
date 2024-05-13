package handlers

import (
	"bytes"
	"io"
	"net/http"
	"time"
	"wannabe/config"
	"wannabe/providers"
	recordEntities "wannabe/record/entities"
	record "wannabe/record/services"

	"github.com/AdguardTeam/gomitmproxy"
)

func WannabeOnResponse(config config.Config, storageProvider providers.StorageProvider) WannabeOnResponseHandler {
	return func(session *gomitmproxy.Session) *http.Response {
		request := session.Request()
		response := session.Response()

		if request.Method == "CONNECT" {
			// FIXME what's right here ???
			return nil
			// return response
		}

		// FIXME block processing based on sessions.Prop
		if _, blocked := session.GetProp("blocked"); blocked {
			return nil
		}
		if _, responseFromRecord := session.GetProp("responseFromRecord"); responseFromRecord {
			return nil
		}

		host := request.URL.Host
		// REVIEW use session.Props ?
		wannabe := config.Wannabes[host]

		// REVIEW use single structure in session.Props to pass everything needed between OnRequest and OnResponse ?
		hash, _ := session.GetProp("hash")
		curl, _ := session.GetProp("curl")

		requestBody, err := io.ReadAll(request.Body)
		if err != nil {
			// return internalError(session, request, err)
			return nil
		}
		defer request.Body.Close()

		// set body back to request
		request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			// return internalError(session, request, err)
			return nil
		}
		// FIXME throws error !!!
		defer response.Body.Close()

		// set body back to response
		response.Body = io.NopCloser(bytes.NewBuffer(responseBody))

		timestamp := time.Now()

		recordPayload := recordEntities.GenerateRecordPayload{
			Hash:            hash.(string),
			Curl:            curl.(string),
			HttpMethod:      request.Method,
			Host:            request.URL.Host,
			Path:            request.URL.Path,
			Query:           request.URL.Query(),
			RequestHeaders:  request.Header,
			RequestBody:     requestBody,
			StatusCode:      response.StatusCode,
			ResponseHeaders: response.Header,
			ResponseBody:    responseBody,
			Timestamp: recordEntities.Timestamp{
				Unix: timestamp.Unix(),
				UTC:  timestamp.UTC(),
			},
		}

		record, err := record.GenerateRecord(wannabe.Records, recordPayload)
		if err != nil {
			// return internalError(session, request, err)
			return nil
		}

		err = storageProvider.InsertRecords([][]byte{record}, []string{hash.(string)}, host)
		if err != nil {
			// return internalError(session, request, err)
			return nil
		}

		// FIXME what's right here ???
		// return nil
		return nil
	}
}
