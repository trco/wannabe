package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Run("load config", func(t *testing.T) {
		configPath, _ := createTestConfigFile("config.json")
		os.Setenv("CONFIG_PATH", configPath)
		os.Setenv("RUNNING_IN_CONTAINER", "true")
		defer os.Unsetenv("CONFIG_PATH")
		defer os.Remove(configPath)
		defer os.Unsetenv("RUNNING_IN_CONTAINER")

		config, _ := LoadConfig()

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
	configPath, _ := createTestConfigFile("config.json")
	os.Setenv("CONFIG_PATH", configPath)
	defer os.Unsetenv("CONFIG_PATH")
	defer os.Remove(configPath)

	tests := []struct {
		name          string
		configPath    string
		defaultConfig Config
		wantConfig    Config
		wantErr       bool
	}{
		{
			name:          "non-existing config file",
			configPath:    "non_existing_config.json",
			defaultConfig: defaultConfig,
			wantConfig:    defaultConfig,
			wantErr:       true,
		},
		{
			name:          "existing config file",
			configPath:    configPath,
			defaultConfig: defaultConfig,
			wantConfig:    wantConfig,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := loadConfigFromFile(tt.configPath, &tt.defaultConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadConfigFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.defaultConfig, tt.wantConfig) {
				t.Errorf("loadConfigFromFile() = %v, want %v", tt.defaultConfig, tt.wantConfig)
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	invalidConfig := Config{
		Mode: "mixed",
		StorageProvider: StorageProvider{
			Type: "filesystem",
			FilesystemConfig: FilesystemConfig{
				Folder:           "records",
				RegenerateFolder: "",
				Format:           "json",
			},
		},
	}

	invalidConfigCustomValidation := *deepCopyConfig(wantConfig)
	wannabe := invalidConfigCustomValidation.Wannabes["testApi"]
	wannabe.RequestMatching = RequestMatching{
		Headers: Headers{
			Include: []string{"Authorization"},
		},
	}
	invalidConfigCustomValidation.Wannabes["testApi"] = wannabe

	tests := []struct {
		name        string
		config      Config
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

	slice := []string{"a", "b", "c"}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := contains(slice, tt.element)

			if got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

var zero = 0
var wantConfig = Config{
	Mode: "mixed",
	StorageProvider: StorageProvider{
		Type: "filesystem",
		FilesystemConfig: FilesystemConfig{
			Folder:           "records",
			RegenerateFolder: "records/regenerated",
			Format:           "json",
		},
	},
	Wannabes: map[string]Wannabe{
		"testApi": {
			RequestMatching: RequestMatching{
				Host: Host{
					Wildcards: []WildcardIndex{
						{Index: &zero, Placeholder: "placeholder"},
					},
				},
				Query: Query{
					Wildcards: []WildcardKey{
						{Key: "status", Placeholder: "placeholder"},
					},
					Regexes: []Regex{
						{Pattern: "app=1", Placeholder: "app=123"},
					},
				},
			},
			Records: Records{
				Headers: HeadersToRecord{
					Exclude: []string{"Authorization"},
				},
			},
		},
	},
}

var defaultConfig = Config{
	Mode: "mixed",
	StorageProvider: StorageProvider{
		Type: "filesystem",
		FilesystemConfig: FilesystemConfig{
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

func deepCopyConfig(c Config) *Config {
	copied := c
	copied.Wannabes = make(map[string]Wannabe, len(c.Wannabes))
	for k, v := range c.Wannabes {
		copied.Wannabes[k] = v
	}
	return &copied
}
