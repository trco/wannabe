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

func (redis RedisProvider) ReadRecords(hashes []string, key string) ([][]byte, error) {
	return [][]byte{}, nil
}

func (redis RedisProvider) InsertRecords(records [][]byte, hashes []string, key string) error {
	return nil
}

func (redis RedisProvider) DeleteRecords(hashes []string, key string) error {
	return nil
}

func (redis RedisProvider) GetHashes(key string) ([]string, error) {
	return []string{}, nil
}
