package storage

import (
	"fmt"

	"github.com/trco/wannabe/internal/config"
)

type Provider interface {
	ReadRecords(subfolder string, hashes []string) ([][]byte, error)
	InsertRecords(subfolder string, hashes []string, records [][]byte, isRegenerate bool) error
	DeleteRecords(subfolder string, hashes []string) error
	GetHashes(subfolder string) ([]string, error)
	GetHostsAndHashes() ([]HostAndHashes, error)
}

type HostAndHashes struct {
	Host      string   `json:"host"`
	HashCount int      `json:"hash_count"`
	Hashes    []string `json:"hashes"`
}

func ProviderFactory(config config.StorageProvider) (Provider, error) {
	if config.Type == "filesystem" {
		storage := FilesystemProvider{
			Config: config,
		}

		return storage, nil
	}

	return nil, fmt.Errorf("generation of '%s' storage provider failed", config.Type)
}
