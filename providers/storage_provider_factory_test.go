package providers

import (
	"reflect"
	"testing"
	"wannabe/config"
)

func TestStorageProviderFactory(t *testing.T) {
	configInvalid := config.StorageProvider{
		Type: "test",
	}

	_, err := StorageProviderFactory(configInvalid)
	if err == nil {
		t.Errorf("Generation of storage provider with invalid config did not throw error.")
	}

	configFilesystem := config.StorageProvider{
		Type: "filesystem",
		FilesystemConfig: config.FilesystemConfig{
			Folder: "records",
			Format: "json",
		},
	}

	storageProvider, _ := StorageProviderFactory(configFilesystem)
	expectedFilesystemProvider := FilesystemProvider{
		Config: configFilesystem,
	}

	if !reflect.DeepEqual(expectedFilesystemProvider, storageProvider) {
		t.Errorf("Expected storage provider: %v, Actual storage provider: %v", expectedFilesystemProvider, storageProvider)
	}

	configRedis := config.StorageProvider{
		Type: "redis",
		RedisConfig: config.RedisConfig{
			Database: "db10",
		},
	}

	storageProvider, _ = StorageProviderFactory(configRedis)
	expectedRedisProvider := RedisProvider{}

	if !reflect.DeepEqual(expectedRedisProvider, storageProvider) {
		t.Errorf("Expected storage provider: %v, Actual storage provider: %v", expectedRedisProvider, storageProvider)
	}
}
