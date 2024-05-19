package providers

import (
	"reflect"
	"testing"
	"wannabe/types"
)

func TestStorageProviderFactory(t *testing.T) {
	invalidConfig := types.Config{
		StorageProvider: types.StorageProvider{
			Type: "test",
		},
	}

	_, err := StorageProviderFactory(invalidConfig)
	if err == nil {
		t.Errorf("generation of storage provider with invalid config did not throw error")
	}

	validConfig := types.Config{
		StorageProvider: types.StorageProvider{
			Type: "filesystem",
			FilesystemConfig: types.FilesystemConfig{
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

	validConfig.StorageProvider = types.StorageProvider{
		Type:       "filesystem",
		Regenerate: true,
		FilesystemConfig: types.FilesystemConfig{
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

	validConfig.StorageProvider = types.StorageProvider{
		Type: "redis",
		RedisConfig: types.RedisConfig{
			Database: "db10",
		},
	}

	storageProvider, _ = StorageProviderFactory(validConfig)
	expectedRedisProvider := RedisProvider{}

	if !reflect.DeepEqual(expectedRedisProvider, storageProvider) {
		t.Errorf("expected storage provider: %v, actual storage provider: %v", expectedRedisProvider, storageProvider)
	}
}
