package providers

import (
	"wannabe/types"
)

type RedisProvider struct {
	Config types.StorageProvider
}

func (redis RedisProvider) GetConfig() types.StorageProvider {
	return redis.Config
}

func (redis RedisProvider) ReadRecords(key string, hashes []string) ([][]byte, error) {
	return [][]byte{}, nil
}

func (redis RedisProvider) InsertRecords(key string, hashes []string, records [][]byte) error {
	return nil
}

func (redis RedisProvider) DeleteRecords(key string, hashes []string) error {
	return nil
}

func (redis RedisProvider) GetHashes(key string) ([]string, error) {
	return []string{}, nil
}

func (fsp RedisProvider) GetHostsAndHashes() ([]HostAndHashes, error) {
	return []HostAndHashes{}, nil
}
