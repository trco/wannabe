package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
	"wannabe/types"
)

func TestLoadConfig(t *testing.T) {
	configFilename, _ := createTestConfigFile("config.json")

	config, _ := LoadConfig(configFilename)

	if !reflect.DeepEqual(testConfig, config) {
		t.Errorf("expected config: %v, actual config: %v", testConfig, config)
	}
}

func TestSetConfigDefaults(t *testing.T) {
	config := setConfigDefaults()

	defaultConfig := types.Config{
		Mode: "mixed",
		StorageProvider: types.StorageProvider{
			Type:       "filesystem",
			Regenerate: false,
			FilesystemConfig: types.FilesystemConfig{
				Folder:           "records",
				RegenerateFolder: "records/regenerated",
				Format:           "json",
			},
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
	configFilename, _ := createTestConfigFile("config.json")
	defer os.Remove(configFilename)

	config, err := loadConfigFromFile(configFilename, configWithDefaults)

	if err != nil {
		t.Errorf("loading of config from existing config.json failed")
	}

	if !reflect.DeepEqual(testConfig, config) {
		t.Errorf("expected config: %v, actual config: %v", testConfig, config)
	}
}

func TestValidateConfig(t *testing.T) {
	// valid config
	err := validateConfig(testConfig)
	if err != nil {
		t.Errorf("valid config validated as invalid")
	}

	// invalid config
	invalidConfig := types.Config{
		StorageProvider: types.StorageProvider{
			Type:       "filesystem",
			Regenerate: false,
			FilesystemConfig: types.FilesystemConfig{
				Folder:           "records",
				RegenerateFolder: "",
				Format:           "json",
			},
		},
	}

	err = validateConfig(invalidConfig)
	if err == nil {
		t.Errorf("invalid config validated as valid")
	}

	// invalid config failing on custom validation
	invalidConfig = testConfig
	wannabe := invalidConfig.Wannabes["testApi"]
	wannabe.RequestMatching = types.RequestMatching{
		Headers: types.Headers{
			Include: []string{"Authorization"},
		},
	}
	invalidConfig.Wannabes["testApi"] = wannabe

	expectedErr := "Key: 'Config.Wannabes' Error:Field validation for 'Wannabes' failed on the 'headers_included_excluded' tag"
	err = validateConfig(invalidConfig)

	if err.Error() != expectedErr {
		t.Errorf("expected error: %s, actual error: %s", expectedErr, err.Error())
	}
}

// reusable variables and methods

var zero = 0
var testConfig = types.Config{
	Mode: "mixed",
	StorageProvider: types.StorageProvider{
		Type:       "filesystem",
		Regenerate: false,
		FilesystemConfig: types.FilesystemConfig{
			Folder:           "records",
			RegenerateFolder: "",
			Format:           "json",
		},
	},
	Wannabes: map[string]types.Wannabe{
		"testApi": {
			RequestMatching: types.RequestMatching{
				Host: types.Host{
					Wildcards: []types.WildcardIndex{
						{Index: &zero, Placeholder: "{placeholder}"},
					},
				},
				Query: types.Query{
					Wildcards: []types.WildcardKey{
						{Key: "status", Placeholder: "{placeholder}"},
					},
					Regexes: []types.Regex{
						{Pattern: "app=1", Placeholder: "app=123"},
					},
				},
			},
			Records: types.Records{
				Headers: types.HeadersToRecord{
					Exclude: []string{"Authorization"},
				},
			},
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
