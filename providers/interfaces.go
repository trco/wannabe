package providers

type StorageProvider interface {
	ReadRecords(subfolder string, hashes []string) ([][]byte, error)
	InsertRecords(subfolder string, hashes []string, records [][]byte) error
	DeleteRecords(subfolder string, hashes []string) error
	GetHashes(subfolder string) ([]string, error)
}
