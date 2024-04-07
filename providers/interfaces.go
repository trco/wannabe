package providers

type StorageProvider interface {
	ReadRecords(hashes []string) ([][]byte, error)
	InsertRecords(hashes []string, records [][]byte) error
	DeleteRecords(hashes []string) error
	GetHashes() ([]string, error)
}
