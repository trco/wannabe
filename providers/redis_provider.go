package providers

import (
	"wannabe/config"
)

type RedisProvider struct {
	Config config.StorageProvider
}

func (redis RedisProvider) GetConfig() config.StorageProvider {
	return redis.Config
}

func (redis RedisProvider) ReadRecords(hashes []string) ([][]byte, error) {
	return [][]byte{}, nil
}

func (redis RedisProvider) InsertRecords(hashes []string, records [][]byte) error {
	return nil
}

func (redis RedisProvider) GetHashes() ([]string, error) {
	return []string{}, nil
}
