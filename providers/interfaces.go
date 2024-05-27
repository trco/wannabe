package providers

type StorageProvider interface {
	ReadRecords(subfolder string, hashes []string) ([][]byte, error)
	InsertRecords(subfolder string, hashes []string, records [][]byte, isRegenerate bool) error
	DeleteRecords(subfolder string, hashes []string) error
	GetHashes(subfolder string) ([]string, error)
	GetHostsAndHashes() ([]HostAndHashes, error)
}

type HostAndHashes struct {
	Host      string   `json:"host"`
	HashCount int      `json:"hash_count"`
	Hashes    []string `json:"hashes"`
}
