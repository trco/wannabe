package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	configFilepath, _ := createTestConfigFile("config.json")

	config, _ := LoadConfig(configFilepath)

	if !reflect.DeepEqual(testConfig, config) {
		t.Errorf("expected config: %v, actual config: %v", testConfig, config)
	}
}

func TestSetConfigDefaults(t *testing.T) {
	config := setConfigDefaults()

	defaultConfig := Config{
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

	if !reflect.DeepEqual(defaultConfig, config) {
		t.Errorf("expected config: %v, actual config: %v", defaultConfig, config)
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	// load config from non-existing config.json
	configWithDefaults := setConfigDefaults()

	_, err := loadConfigFromFile("config.json", configWithDefaults)
	if err == nil {
		t.Errorf("loading of config from non-existing file did not throw error")
	}

	// load config from existing config.json
	configFilepath, _ := createTestConfigFile("config.json")
	defer os.Remove(configFilepath)

	config, err := loadConfigFromFile(configFilepath, configWithDefaults)

	if err != nil {
		t.Errorf("loading of config from existing config.json failed")
	}

	if !reflect.DeepEqual(testConfig, config) {
		t.Errorf("expected config: %v, actual config: %v", testConfig, config)
	}
}

func TestValidateConfig(t *testing.T) {
	// invalid config
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
		t.Errorf("invalid config validated as valid")
	}

	// invalid config failing on custom validation
	invalidConfig = testConfig
	invalidConfig.RequestMatching = RequestMatching{
		Headers: Headers{
			Include: []string{"Authorization"},
		},
	}

	expectedErr := "Key: 'Config.RequestMatching.Headers.Include' Error:Field validation for 'Include' failed on the 'the_same_header_defined_in_records_headers_exclude' tag"
	err = validateConfig(invalidConfig)

	if err.Error() != expectedErr {
		t.Errorf("expected error: %s, actual error: %s", expectedErr, err.Error())
	}

	// valid config
	err = validateConfig(testConfig)
	if err != nil {
		t.Errorf("valid config validated as invalid")
	}

}

// reusable variables and methods

var zero = 0
var testConfig = Config{
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
	Records: Records{
		Headers: HeadersToExclude{
			Exclude: []string{"Authorization"},
		},
	},
}

func createTestConfigFile(filename string) (string, error) {
	tempFile, err := os.CreateTemp("", filename)
	if err != nil {
		fmt.Println("error creating temporary config.json file:", err)
		return "", err
	}
	defer tempFile.Close()

	jsonData, err := json.Marshal(testConfig)
	if err != nil {
		fmt.Println("error encoding config to json:", err)
		return "", err
	}

	err = os.WriteFile(tempFile.Name(), jsonData, 0644)
	if err != nil {
		fmt.Println("error writing config json to temporary config.json file:", err)
		return "", nil
	}

	path := tempFile.Name()

	return path, nil
}
