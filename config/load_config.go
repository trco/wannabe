package config

import (
	"fmt"
	"wannabe/utils"

	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
)

func LoadConfig(configFilename string) (Config, error) {
	config := setConfigDefaults()

	config, err := loadConfigFromFile(configFilename, config)
	if err != nil {
		return Config{}, fmt.Errorf("failed loading config from file: %v", err)
	}

	err = validateConfig(config)
	if err != nil {
		return Config{}, fmt.Errorf("failed validating config: %v", err)
	}

	return config, nil
}

func setConfigDefaults() Config {
	return Config{
		Mode: "mixed",
		StorageProvider: StorageProvider{
			Type:       "filesystem",
			Regenerate: false,
			FilesystemConfig: FilesystemConfig{
				Folder:           "records",
				RegenerateFolder: "records/regenerated",
				Format:           "json",
			},
		},
	}
}

func loadConfigFromFile(configFilename string, config Config) (Config, error) {
	f := file.Provider(configFilename)
	var k = koanf.New(".")

	loadConfig := func() error {
		// overwrites config defaults
		err := k.Load(f, json.Parser())
		if err != nil {
			return err
		}

		err = k.Unmarshal("", &config)
		if err != nil {
			return err
		}

		return nil
	}

	err := loadConfig()
	if err != nil {
		return Config{}, err
	}

	k.Print()

	return config, nil
}

var validate *validator.Validate

func validateConfig(config Config) error {
	validate = validator.New()

	validate.RegisterValidation("headers_included_excluded", validateWannabeHeadersConfig)

	err := validate.Struct(config)
	if err != nil {
		return err
	}

	return nil
}

func validateWannabeHeadersConfig(fl validator.FieldLevel) bool {
	wannabes := fl.Field().Interface().(map[string]Wannabe)

	for _, wannabe := range wannabes {
		headersInclude := wannabe.RequestMatching.Headers.Include
		headersExclude := wannabe.Records.Headers.Exclude

		for _, i := range headersInclude {
			if utils.Contains(headersExclude, i) {
				return false

			}
		}
	}

	return true
}
