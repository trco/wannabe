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
		storageProvider := FilesystemProvider{
			Config: spc,
		}

		err := storageProvider.CreateFolders()
		if err != nil {
			return nil, fmt.Errorf("StorageProviderFactory: failed creating folders: %v", err)
		}

		return storageProvider, nil
	}

	if spc.Type == "redis" {
		return RedisProvider{}, nil
	}

	return nil, &StorageProviderGenerationError{
		Type: spc.Type,
	}
}
