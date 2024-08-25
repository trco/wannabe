package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
)

func LoadConfig() (Config, error) {
	config := setConfigDefaults()

	configPath, err := getConfigPath()
	if err != nil {
		return Config{}, fmt.Errorf("failed getting config path: %v", err)
	}

	if configPath != "" {
		err := loadConfigFromFile(configPath, &config)
		if err != nil {
			return Config{}, fmt.Errorf("failed loading config from file: %v", err)
		}
	}

	err = validateConfig(config)
	if err != nil {
		return Config{}, fmt.Errorf("failed validating config: %v", err)
	}

	return config, nil
}

func getConfigPath() (string, error) {
	var configPath string

	if os.Getenv(RunningInContainer) == "" {
		// check if config.json exists
		_, err := os.Stat("config.json")
		if err != nil && !os.IsNotExist(err) {
			return "", fmt.Errorf("failed checking if config.json file exists in the root folder")
		} else if os.IsNotExist(err) {
			return "", nil
		}

		return "config.json", nil
	}

	configPath = os.Getenv(ConfigPath)
	if configPath == "" {
		return "", fmt.Errorf("%v env variable not set", ConfigPath)
	}

	return configPath, nil
}

func setConfigDefaults() Config {
	return Config{
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
}

func loadConfigFromFile(configPath string, config *Config) error {
	f := file.Provider(configPath)
	var k = koanf.New(".")

	loadConfig := func() error {
		// overwrites config defaults
		err := k.Load(f, json.Parser())
		if err != nil {
			return err
		}

		err = k.Unmarshal("", config)
		if err != nil {
			return err
		}

		return nil
	}

	err := loadConfig()
	if err != nil {
		return err
	}

	return nil
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
			if contains(headersExclude, i) {
				return false

			}
		}
	}

	return true
}

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

const (
	RunningInContainer = "RUNNING_IN_CONTAINER"
	ConfigPath         = "CONFIG_PATH"
	CertPath           = "CERT_PATH"
	CertKeyPath        = "CERT_KEY_PATH"
)

type Config struct {
	Mode            string             `koanf:"mode" validate:"required,oneof=proxy server mixed"`
	StorageProvider StorageProvider    `koanf:"storageProvider" validate:"required"`
	Wannabes        map[string]Wannabe `koanf:"wannabes" validate:"headers_included_excluded,dive"`
}

const (
	ServerMode = "server"
	MixedMode  = "mixed"
	ProxyMode  = "proxy"
)

type StorageProvider struct {
	Type             string           `koanf:"type" validate:"required,oneof=filesystem"`
	FilesystemConfig FilesystemConfig `koanf:"filesystemConfig" validate:"required_if=Type filesystem,omitempty"`
}

type FilesystemConfig struct {
	Folder           string `koanf:"folder" validate:"required"`
	RegenerateFolder string `koanf:"regenerateFolder"`
	Format           string `koanf:"format" validate:"required,oneof=json"`
}

type Wannabe struct {
	RequestMatching RequestMatching `koanf:"requestMatching"`
	Records         Records         `koanf:"records"`
}

type RequestMatching struct {
	Host    Host    `koanf:"host"`
	Path    Path    `koanf:"path"`
	Query   Query   `koanf:"query"`
	Body    Body    `koanf:"body"`
	Headers Headers `koanf:"headers"`
}

type Host struct {
	Wildcards []WildcardIndex `koanf:"wildcards" validate:"gte=0,dive"`
	Regexes   []Regex         `koanf:"regexes" validate:"gte=0,dive"`
}

type Path struct {
	Wildcards []WildcardIndex `koanf:"wildcards" validate:"gte=0,dive"`
	Regexes   []Regex         `koanf:"regexes" validate:"gte=0,dive"`
}

type Query struct {
	Wildcards []WildcardKey `koanf:"wildcards" validate:"gte=0,dive"`
	Regexes   []Regex       `koanf:"regexes" validate:"gte=0,dive"`
}

type Body struct {
	Regexes []Regex `koanf:"regexes" validate:"gte=0,dive"`
}

type Headers struct {
	Include   []string      `koanf:"include" validate:"gte=0,dive"`
	Wildcards []WildcardKey `koanf:"wildcards" validate:"gte=0,dive"`
}

type WildcardIndex struct {
	// pointer is used to pass "required" validation with "0" value
	Index       *int   `koanf:"index" validate:"required,numeric,min=0"`
	Placeholder string `koanf:"placeholder" validate:"ascii"`
}

type WildcardKey struct {
	Key         string `koanf:"key" validate:"required,alphanum"`
	Placeholder string `koanf:"placeholder" validate:"ascii"`
}

type WildcardPath struct {
	Path        string `koanf:"path" validate:"required,uri"`
	Placeholder string `koanf:"placeholder" validate:"ascii"`
}

type Regex struct {
	Pattern     string `koanf:"pattern" validate:"required,ascii"`
	Placeholder string `koanf:"placeholder" validate:"ascii"`
}

type Records struct {
	Headers HeadersToRecord `koanf:"headers"`
}

type HeadersToRecord struct {
	Exclude []string `koanf:"exclude"`
}
