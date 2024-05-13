package providers

import (
	"reflect"
	"testing"
	"wannabe/config"
)

func TestStorageProviderFactory(t *testing.T) {
	invalidConfig := config.Config{
		StorageProvider: config.StorageProvider{
			Type: "test",
		},
	}

	_, err := StorageProviderFactory(invalidConfig)
	if err == nil {
		t.Errorf("generation of storage provider with invalid config did not throw error")
	}

	validConfig := config.Config{
		StorageProvider: config.StorageProvider{
			Type: "filesystem",
			FilesystemConfig: config.FilesystemConfig{
				Folder:           "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T",
				RegenerateFolder: "",
				Format:           "json",
			},
		},
	}

	storageProvider, _ := StorageProviderFactory(validConfig)
	expectedFilesystemProvider := FilesystemProvider{
		Config: validConfig,
	}

	if !reflect.DeepEqual(expectedFilesystemProvider, storageProvider) {
		t.Errorf("expected storage provider: %v, actual storage provider: %v", expectedFilesystemProvider, storageProvider)
	}

	validConfig.StorageProvider = config.StorageProvider{
		Type:       "filesystem",
		Regenerate: true,
		FilesystemConfig: config.FilesystemConfig{
			Folder:           "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T",
			RegenerateFolder: "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T/regenerate",
			Format:           "json",
		},
	}

	storageProvider, _ = StorageProviderFactory(validConfig)
	expectedRegenerateProvider := FilesystemProvider{
		Config: validConfig,
	}

	if !reflect.DeepEqual(expectedRegenerateProvider, storageProvider) {
		t.Errorf("expected storage provider: %v, actual storage provider: %v", expectedRegenerateProvider, storageProvider)
	}

	validConfig.StorageProvider = config.StorageProvider{
		Type: "redis",
		RedisConfig: config.RedisConfig{
			Database: "db10",
		},
	}

	storageProvider, _ = StorageProviderFactory(validConfig)
	expectedRedisProvider := RedisProvider{}

	if !reflect.DeepEqual(expectedRedisProvider, storageProvider) {
		t.Errorf("expected storage provider: %v, actual storage provider: %v", expectedRedisProvider, storageProvider)
	}
}
