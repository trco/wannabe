package config

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
)

func Load(configFilepath string) Config {
	config := setConfigDefaults()

	config, err := loadConfigFromFile(configFilepath, config)
	if err != nil {
		log.Fatalf("fatal error loading config: %v", err)
	}

	err = validateConfig(config)
	if err != nil {
		log.Fatalf("fatal error validating config: %v", err)
	}

	return config
}

func setConfigDefaults() Config {
	return Config{
		StorageProvider: StorageProvider{
			Type: "filesystem",
			FilesystemConfig: FilesystemConfig{
				Folder: "records",
				Format: "json",
			},
		},
		Read: Read{
			Enabled:     true,
			FailOnError: false,
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

	f.Watch(func(event interface{}, err error) {
		if err != nil {
			log.Fatalf("fatal error watching config file: %v", err)
		}

		err = loadConfig()
		if err != nil {
			log.Fatalf("fatal error updating config: %v", err)
		}

		err = validateConfig(config)
		if err != nil {
			log.Fatalf("fatal error validating updated config: %v", err)
		}

		log.Println("config successfully updated")
	})

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
