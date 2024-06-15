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

	t.Run("load config", func(t *testing.T) {
		config, _ := LoadConfig(configFilename)

		if !reflect.DeepEqual(config, wantConfig) {
			t.Errorf("LoadConfig() = %v, want %v", config, wantConfig)
		}
	})
}

func TestSetConfigDefaults(t *testing.T) {
	t.Run("default config", func(t *testing.T) {
		config := setConfigDefaults()

		if !reflect.DeepEqual(config, defaultConfig) {
			t.Errorf("setConfigDefaults() = %v, want %v", config, defaultConfig)
		}
	})
}

func TestLoadConfigFromFile(t *testing.T) {
	configFilename, _ := createTestConfigFile("config.json")
	defer os.Remove(configFilename)

	tests := []struct {
		name          string
		filename      string
		defaultConfig types.Config
		wantConfig    types.Config
		wantErr       bool
	}{
		{
			name:          "non-existing config file",
			filename:      "non_existing_config.json",
			defaultConfig: defaultConfig,
			wantConfig:    types.Config{},
			wantErr:       true,
		},
		{
			name:          "existing config file",
			filename:      configFilename,
			defaultConfig: defaultConfig,
			wantConfig:    wantConfig,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadConfigFromFile(tt.filename, tt.defaultConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadConfigFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantConfig) {
				t.Errorf("loadConfigFromFile() = %v, want %v", got, tt.wantConfig)
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	invalidConfig := types.Config{
		Mode: "mixed",
		StorageProvider: types.StorageProvider{
			Type: "filesystem",
			FilesystemConfig: types.FilesystemConfig{
				Folder:           "records",
				RegenerateFolder: "",
				Format:           "json",
			},
		},
	}

	invalidConfigCustomValidation := *deepCopyConfig(wantConfig)
	wannabe := invalidConfigCustomValidation.Wannabes["testApi"]
	wannabe.RequestMatching = types.RequestMatching{
		Headers: types.Headers{
			Include: []string{"Authorization"},
		},
	}
	invalidConfigCustomValidation.Wannabes["testApi"] = wannabe

	tests := []struct {
		name        string
		config      types.Config
		expectedErr string
	}{
		{
			name:        "invalid config",
			config:      invalidConfig,
			expectedErr: "Key: 'Config.StorageProvider.FilesystemConfig.RegenerateFolder' Error:Field validation for 'RegenerateFolder' failed on the 'required' tag",
		},
		{
			name:        "valid config",
			config:      wantConfig,
			expectedErr: "",
		},
		{
			name:        "invalid config failing on custom validation",
			config:      invalidConfigCustomValidation,
			expectedErr: "Key: 'Config.Wannabes' Error:Field validation for 'Wannabes' failed on the 'headers_included_excluded' tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)
			if err != nil && err.Error() != tt.expectedErr {
				t.Errorf("validateConfig() = %v, want %v", err.Error(), tt.expectedErr)
			}
		})
	}
}

func TestContains(t *testing.T) {
	slice := []string{"a", "b", "c"}

	tests := []struct {
		name    string
		element string
		want    bool
	}{
		{
			name:    "element exists in slice",
			element: "b",
			want:    true,
		},
		{
			name:    "element does not exist in slice",
			element: "d",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := contains(slice, tt.element)

			if got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

// reusable variables and methods
var zero = 0
var wantConfig = types.Config{
	Mode: "mixed",
	StorageProvider: types.StorageProvider{
		Type: "filesystem",
		FilesystemConfig: types.FilesystemConfig{
			Folder:           "records",
			RegenerateFolder: "records/regenerated",
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

var defaultConfig = types.Config{
	Mode: "mixed",
	StorageProvider: types.StorageProvider{
		Type: "filesystem",
		FilesystemConfig: types.FilesystemConfig{
			Folder:           "records",
			RegenerateFolder: "records/regenerated",
			Format:           "json",
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

	jsonData, err := json.Marshal(wantConfig)
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

func deepCopyConfig(c types.Config) *types.Config {
	copied := c
	copied.Wannabes = make(map[string]types.Wannabe, len(c.Wannabes))
	for k, v := range c.Wannabes {
		copied.Wannabes[k] = v
	}
	return &copied
}
