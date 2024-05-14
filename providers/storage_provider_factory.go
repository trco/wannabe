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

func StorageProviderFactory(config config.Config) (StorageProvider, error) {
	if config.StorageProvider.Type == "filesystem" {
		storageProvider := FilesystemProvider{
			Config: config,
		}

		return storageProvider, nil
	}

	if config.StorageProvider.Type == "redis" {
		return RedisProvider{}, nil
	}

	return nil, &StorageProviderGenerationError{
		Type: config.StorageProvider.Type,
	}
}
