package providers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"wannabe/config"
)

type FilesystemProvider struct {
	Config config.StorageProvider
}

func (fsp FilesystemProvider) CreateFolders() error {
	err := os.Mkdir(fsp.Config.FilesystemConfig.Folder, 0750)
	if err != nil && !os.IsExist(err) {
		return err
	}

	if fsp.Config.Regenerate {
		err = os.Mkdir(fsp.Config.FilesystemConfig.RegenerateFolder, 0750)
		if err != nil && !os.IsExist(err) {
			return filesystemProviderErr("failed creating regenerate folder or missing 'storageProvider.filesystemConfig.regenerateFolder' config", err)
		}
	}

	return nil
}

// TODO ? bulk read using goroutines and channels
func (fsp FilesystemProvider) ReadRecords(hashes []string) ([][]byte, error) {
	var records [][]byte

	// TODO all or nothing ?
	for _, hash := range hashes {
		filepath := fsp.generateFilepath(hash)

		_, err := os.Stat(filepath)
		if err != nil {
			return nil, filesystemProviderErr("failed checking if file exists", err)
		}

		file, err := os.Open(filepath)
		if err != nil {
			return nil, filesystemProviderErr("failed opening file", err)
		}
		defer file.Close()

		record, err := io.ReadAll(file)
		if err != nil {
			return nil, filesystemProviderErr("failed reading file", err)
		}

		records = append(records, record)
	}

	return records, nil
}

// TODO ? bulk insert using goroutines and channels
func (fsp FilesystemProvider) InsertRecords(hashes []string, records [][]byte) error {
	// TODO all or nothing ?
	for index, record := range records {
		filepath := fsp.generateFilepath(hashes[index])
		if fsp.Config.Regenerate {
			filepath = fsp.generateFilepathRegenerate(hashes[index])
		}

		_, err := os.Create(filepath)
		if err != nil {
			return filesystemProviderErr("failed creating file", err)
		}

		err = os.WriteFile(filepath, record, 0644)
		if err != nil {
			return filesystemProviderErr("failed writing file", err)
		}
	}

	return nil
}

// TODO ? bulk delete using goroutines and channels
func (fsp FilesystemProvider) DeleteRecords(hashes []string) error {
	for _, hash := range hashes {
		filepath := fsp.generateFilepath(hash)

		err := os.Remove(filepath)
		if err != nil {
			return filesystemProviderErr("failed deleting file", err)
		}
	}

	return nil
}

func (fsp FilesystemProvider) GetHashes() ([]string, error) {
	folder := fsp.Config.FilesystemConfig.Folder

	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, filesystemProviderErr("failed reading folder", err)
	}

	var hashes []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// get the file name without extension
		filename := file.Name()
		hash := filename[:len(filename)-len(filepath.Ext(filename))]

		hashes = append(hashes, hash)
	}

	return hashes, nil
}

func (fsp FilesystemProvider) generateFilepath(hash string) string {
	folder := fsp.Config.FilesystemConfig.Folder
	format := fsp.Config.FilesystemConfig.Format

	return folder + "/" + hash + "." + format
}

func (fsp FilesystemProvider) generateFilepathRegenerate(hash string) string {
	folder := fsp.Config.FilesystemConfig.RegenerateFolder
	format := fsp.Config.FilesystemConfig.Format

	return folder + "/" + hash + "." + format
}

func filesystemProviderErr(message string, err error) error {
	return fmt.Errorf("FileSystemProvider: %s: %v", message, err)
}
