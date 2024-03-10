package config

type Config struct {
	// // pointers are used for "required_without_all" validation to work on structs
	// FilesystemStorage *FilesystemStorage `koanf:"filesystemStorage" validate:"required_without_all=RedisStorage"`
	// RedisStorage      *RedisStorage      `koanf:"redis" validate:"required_without_all=FilesystemStorage"`
	StorageProvider StorageProvider `koanf:"storageProvider" validate:"required"`
	// REVIEW add Read as default, does it get overriden when config is loaded from file
	Read            Read            `koanf:"read" validate:"required"`
	Server          string          `koanf:"server" validate:"required,http_url"`
	RequestMatching RequestMatching `koanf:"requestMatching"`
	Records         Records         `koanf:"records"`
	Logger          Logger          `koanf:"logger"`
}

type StorageProvider struct {
	Type             string           `koanf:"type" validate:"required,oneof=filesystem redis"`
	FilesystemConfig FilesystemConfig `koanf:"filesystemConfig" validate:"required_if=Type filesystem,omitempty"`
	RedisConfig      RedisConfig      `koanf:"redisConfig" validate:"required_if=Type redis,omitempty"`
	// Folder   string `koanf:"folder" validate:"required_if=Type filesystem"`
	// Format   string `koanf:"format" validate:"required_if=Type filesystem,oneof=json"`
	// Database string `koanf:"database" validate:"required_if=Type redis"`
}

type FilesystemConfig struct {
	// REVIEW validate:"dirpath" demands existing folder
	Folder string `koanf:"folder" validate:"required"`
	Format string `koanf:"format" validate:"required,oneof=json"`
}

type RedisConfig struct {
	Database string `koanf:"database" validate:"required"`
}

// type FilesystemStorage struct {
// 	// REVIEW validate:"dirpath" demands existing folder
// 	Folder string `koanf:"folder" validate:"required"`
// 	Format string `koanf:"format" validate:"required,oneof=json"`
// }

// type RedisStorage struct {
// 	// REVIEW do allowed database names consist of alphanum characters only ?
// 	Database int `koanf:"db" validate:"required"`
// }

type Read struct {
	// no need to validate as "required" since the default boolean values are set
	Enabled bool `koanf:"enabled" validate:"boolean"`
	// REVIEW prepare command ??? that will validate config and return warning in case Read.Enable: true and Read.FailOnError: true
	// when Read.FailOnError: true the app returns 500 internal error:
	// - if record is not found in storage
	// - if file can't be opened
	// - if file can't be read
	// - if record can't be unmarshaled in PrepareResponse
	FailOnError bool `koanf:"failOnError" validate:"boolean"`
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
	Wildcards []WildcardPath `koanf:"wildcards" validate:"gte=0,dive"`
	Regexes   []Regex        `koanf:"regexes" validate:"gte=0,dive"`
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
	Headers HeadersToExclude `koanf:"headers"`
}

type HeadersToExclude struct {
	Exclude []string `koanf:"exclude"`
}

type Logger struct {
	Enabled  bool   `koanf:"enabled" validate:"boolean"`
	Filepath string `koanf:"filepath"`
	Format   string `koanf:"format"`
}
