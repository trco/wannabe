package providers

import (
	"fmt"
	"wannabe/types"
)

type StorageProviderGenerationError struct {
	Type string
}

func (e *StorageProviderGenerationError) Error() string {
	return fmt.Sprintf("Generation of '%v' storage provider failed.", e.Type)
}

func StorageProviderFactory(config types.Config) (StorageProvider, error) {
	if config.StorageProvider.Type == "filesystem" {
		storageProvider := FilesystemProvider{
			Config: config,
		}

		return storageProvider, nil
	}

	return nil, &StorageProviderGenerationError{
		Type: config.StorageProvider.Type,
	}
}
