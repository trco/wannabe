package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
)

func LoadConfig(configFilepath string) (Config, error) {
	config := setConfigDefaults()

	config, err := loadConfigFromFile(configFilepath, config)
	if err != nil {
		return Config{}, fmt.Errorf("failed loading config: %v", err)
	}

	err = validateConfig(config)
	if err != nil {
		return Config{}, fmt.Errorf("failed validating config: %v", err)
	}

	return config, nil
}

func setConfigDefaults() Config {
	return Config{
		Mode:            "mixed",
		FailOnReadError: false,
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

func loadConfigFromFile(configFilepath string, config Config) (Config, error) {
	f := file.Provider(configFilepath)
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

	validate.RegisterValidation("the_same_header_defined_in_records_headers_exclude", validateHeadersConfig)

	err := validate.Struct(config)
	if err != nil {
		return err
	}

	return nil
}

// custom validation functions
func validateHeadersConfig(fl validator.FieldLevel) bool {
	fieldInclude := fl.Parent().FieldByName(fl.StructFieldName())
	include := fieldInclude.Interface().([]string)

	fieldExclude := fl.Top().FieldByName("Records").FieldByName("Headers").FieldByName("Exclude")
	exclude := fieldExclude.Interface().([]string)

	for _, i := range include {
		for _, e := range exclude {
			if i == e {
				return false
			}
		}
	}

	return true
}
