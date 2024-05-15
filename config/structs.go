package config

type Config struct {
	Mode            string             `koanf:"mode" validate:"required,oneof=proxy server mixed"`
	StorageProvider StorageProvider    `koanf:"storageProvider" validate:"required"`
	Wannabes        map[string]Wannabe `koanf:"wannabes" validate:"required,headers_included_excluded,dive"`
}

const (
	ServerMode = "server"
	MixedMode  = "mixed"
	ProxyMode  = "proxy"
)

type StorageProvider struct {
	Type             string           `koanf:"type" validate:"required,oneof=filesystem redis"`
	Regenerate       bool             `koanf:"regenerate"`
	FilesystemConfig FilesystemConfig `koanf:"filesystemConfig" validate:"required_if=Type filesystem,omitempty"`
	RedisConfig      RedisConfig      `koanf:"redisConfig" validate:"required_if=Type redis,omitempty"`
}

type FilesystemConfig struct {
	Folder           string `koanf:"folder" validate:"required"`
	RegenerateFolder string `koanf:"regenerateFolder"`
	Format           string `koanf:"format" validate:"required,oneof=json"`
}

type RedisConfig struct {
	Database string `koanf:"database" validate:"required"`
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
	Headers HeadersToExclude `koanf:"headers"`
}

type HeadersToExclude struct {
	Exclude []string `koanf:"exclude"`
}
