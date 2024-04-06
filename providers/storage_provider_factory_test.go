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
		t.Errorf("generation of storage provider with invalid config did not throw error")
	}

	configFilesystem := config.StorageProvider{
		Type: "filesystem",
		FilesystemConfig: config.FilesystemConfig{
			Folder:           "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T",
			RegenerateFolder: "",
			Format:           "json",
		},
	}

	storageProvider, _ := StorageProviderFactory(configFilesystem)
	expectedFilesystemProvider := FilesystemProvider{
		Config: configFilesystem,
	}

	if !reflect.DeepEqual(expectedFilesystemProvider, storageProvider) {
		t.Errorf("expected storage provider: %v, actual storage provider: %v", expectedFilesystemProvider, storageProvider)
	}

	configRegenerate := config.StorageProvider{
		Type:       "filesystem",
		Regenerate: true,
		FilesystemConfig: config.FilesystemConfig{
			Folder:           "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T",
			RegenerateFolder: "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T/regenerate",
			Format:           "json",
		},
	}

	storageProvider, _ = StorageProviderFactory(configRegenerate)
	expectedRegenerateProvider := FilesystemProvider{
		Config: configRegenerate,
	}

	if !reflect.DeepEqual(expectedRegenerateProvider, storageProvider) {
		t.Errorf("expected storage provider: %v, actual storage provider: %v", expectedRegenerateProvider, storageProvider)
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
		t.Errorf("expected storage provider: %v, actual storage provider: %v", expectedRedisProvider, storageProvider)
	}
}
