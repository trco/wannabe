package handlers

import "github.com/gofiber/fiber/v2"

type WannabeHandler func(ctx *fiber.Ctx) error

type InternalError struct {
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
}

type RecordProcessingDetails struct {
	Hash    string `json:"hash"`
	Message string `json:"message"`
}

type PostRecordsResponse struct {
	InsertedRecordsCount    int                       `json:"insertedRecordsCount"`
	NotInsertedRecordsCount int                       `json:"notInsertedRecordsCount"`
	RecordProcessingDetails []RecordProcessingDetails `json:"recordProcessingDetails"`
}

type DeleteRecordsResponse struct {
	Message string   `json:"message"`
	Hashes  []string `json:"hashes"`
}

type RegenerateResponse struct {
	Message           string   `json:"message"`
	RegeneratedHashes []string `json:"regeneratedHashes"`
	FailedHashes      []string `json:"failedHashes"`
}
