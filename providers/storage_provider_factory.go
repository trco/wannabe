package providers

import (
	"fmt"
	"wannabe/config"
)

type StorageProviderGenerationError struct {
	Type string
}

func (e *StorageProviderGenerationError) Error() string {
	return fmt.Sprintf("Generation of '%v' storage provider failed.", e.Type)
}

func StorageProviderFactory(spc config.StorageProvider) (StorageProvider, error) {
	if spc.Type == "filesystem" {
		return FilesystemProvider{
			Config: spc,
		}, nil
	}

	if spc.Type == "redis" {
		return RedisProvider{}, nil
	}

	return nil, &StorageProviderGenerationError{
		Type: spc.Type,
	}
}
