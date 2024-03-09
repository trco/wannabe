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

	// var config Config

	config, err := loadConfigFromFile(configFilepath, config)
	if err != nil {
		log.Fatalf("fatal error when loading config from file: %v", err)
	}

	err = validateConfig(config)
	if err != nil {
		log.Fatalf("fatal error when validating config loaded from file: %v", err)
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
		// Logger: Logger{
		// 	Filepath: "./wannabe.log",
		// 	Format:   "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\nResponse body\n${resBody}",
		// },
	}
}

func loadConfigFromFile(configFilepath string, config Config) (Config, error) {
	var k = koanf.New(".")

	// overwrites defaults
	err := k.Load(file.Provider(configFilepath), json.Parser())
	if err != nil {
		return config, err
	}

	k.Unmarshal("", &config)

	return config, nil
}

func validateConfig(config Config) error {
	err := validator.New().Struct(config)
	if err != nil {
		return err
	}

	return nil
}
