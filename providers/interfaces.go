package providers

type StorageProvider interface {
	ReadRecords(hashes []string, host string) ([][]byte, error)
	InsertRecords(records [][]byte, hashes []string, host string) error
	DeleteRecords(hashes []string, host string) error
	GetHashes(host string) ([]string, error)
}
