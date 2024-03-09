package providers

import (
	"wannabe/config"
)

type StorageProvider interface {
	GetConfig() config.StorageProvider
	ReadRecords(hashes []string) ([][]byte, error)
	InsertRecords(hashes []string, records [][]byte) error
	GetHashes() ([]string, error)
}
