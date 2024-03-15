package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func testConfig() Config {
	zero := 0
	return Config{
		StorageProvider: StorageProvider{
			Type:       "filesystem",
			Regenerate: false,
			FilesystemConfig: FilesystemConfig{
				Folder:           "records",
				RegenerateFolder: "",
				Format:           "json",
			},
		},
		Read: Read{
			Enabled:     true,
			FailOnError: true,
		},
		Server: "https://analyticsdata.googleapis.com",
		RequestMatching: RequestMatching{
			Host: Host{
				Wildcards: []WildcardIndex{
					{Index: &zero, Placeholder: "{placeholder}"},
				},
			},
			Query: Query{
				Wildcards: []WildcardKey{
					{Key: "status", Placeholder: "{placeholder}"},
				},
				Regexes: []Regex{
					{Pattern: "app=1", Placeholder: "app=123"},
				},
			},
		},
	}
}

func TestLoad(t *testing.T) {
	configFilepath, _ := createTestConfigFile("config.json")

	config := Load(configFilepath)

	expectedConfig := testConfig()

	if !reflect.DeepEqual(expectedConfig, config) {
		t.Errorf("Expected config: %v, Actual config: %v", expectedConfig, config)
	}
}

func TestSetConfigDefaults(t *testing.T) {
	config := setConfigDefaults()

	expectedConfig := Config{
		StorageProvider: StorageProvider{
			Type:       "filesystem",
			Regenerate: false,
			FilesystemConfig: FilesystemConfig{
				Folder:           "records",
				RegenerateFolder: "",
				Format:           "json",
			},
		},
		Read: Read{
			Enabled:     true,
			FailOnError: false,
		},
	}

	if !reflect.DeepEqual(expectedConfig, config) {
		t.Errorf("Expected config: %v, Actual config: %v", expectedConfig, config)
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	// load config from non-existing config.json
	configWithDefaults := setConfigDefaults()
	_, err := loadConfigFromFile("config.json", configWithDefaults)

	if err == nil {
		t.Errorf("Loading of config from non-existing file did not throw error.")
	}

	configFilepath, _ := createTestConfigFile("config.json")

	// load config from existing config.json
	config, err := loadConfigFromFile(configFilepath, configWithDefaults)

	expectedConfig := testConfig()

	if err != nil {
		t.Errorf("Loading of config from existing config.json failed.")
	}

	if !reflect.DeepEqual(expectedConfig, config) {
		t.Errorf("Expected config: %v, Actual config: %v", expectedConfig, config)
	}
}

func TestValidateConfig(t *testing.T) {
	invalidConfig := Config{
		StorageProvider: StorageProvider{
			Type:       "filesystem",
			Regenerate: false,
			FilesystemConfig: FilesystemConfig{
				Folder:           "records",
				RegenerateFolder: "",
				Format:           "json",
			},
		},
	}

	err := validateConfig(invalidConfig)
	if err == nil {
		t.Errorf("Invalid config validated as valid.")
	}

	validConfig := testConfig()

	err = validateConfig(validConfig)
	if err != nil {
		t.Errorf("Valid config validated as invalid.")
	}
}

func createTestConfigFile(filename string) (string, error) {
	tempFile, err := os.CreateTemp("", filename)
	if err != nil {
		fmt.Println("Error creating temporary config.json file:", err)
		return "", err
	}
	defer tempFile.Close()

	zero := 0
	config := Config{
		StorageProvider: StorageProvider{
			Type:       "filesystem",
			Regenerate: false,
			FilesystemConfig: FilesystemConfig{
				Folder:           "records",
				RegenerateFolder: "",
				Format:           "json",
			},
		},
		Read: Read{
			Enabled:     true,
			FailOnError: true,
		},
		Server: "https://analyticsdata.googleapis.com",
		RequestMatching: RequestMatching{
			Host: Host{
				Wildcards: []WildcardIndex{
					{Index: &zero, Placeholder: "{placeholder}"},
				},
			},
			Query: Query{
				Wildcards: []WildcardKey{
					{Key: "status", Placeholder: "{placeholder}"},
				},
				Regexes: []Regex{
					{Pattern: "app=1", Placeholder: "app=123"},
				},
			},
		},
	}

	jsonData, err := json.Marshal(config)
	if err != nil {
		fmt.Println("Error encoding config to json:", err)
		return "", err
	}

	err = os.WriteFile(tempFile.Name(), jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing config json to temporary config.json file:", err)
		return "", nil
	}

	path := tempFile.Name()

	return path, nil
}
