package providers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"wannabe/config"
)

type FilesystemProvider struct {
	Config config.Config
}

// TODO ? bulk read using goroutines and channels
func (fsp FilesystemProvider) ReadRecords(hashes []string, subfolder string) ([][]byte, error) {
	var records [][]byte

	// TODO all or nothing ?
	for _, hash := range hashes {
		filepath := fsp.generateFilepath(subfolder, hash, false)

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
func (fsp FilesystemProvider) InsertRecords(records [][]byte, hashes []string, subfolder string) error {
	regenerate := fsp.Config.StorageProvider.Regenerate

	err := fsp.createFolder(subfolder, regenerate)
	if err != nil {
		return filesystemProviderErr("failed creating folder", err)
	}

	// TODO all or nothing ?
	for index, record := range records {
		filepath := fsp.generateFilepath(subfolder, hashes[index], regenerate)

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
func (fsp FilesystemProvider) DeleteRecords(hashes []string, subfolder string) error {
	for _, hash := range hashes {
		filepath := fsp.generateFilepath(subfolder, hash, false)

		err := os.Remove(filepath)
		if err != nil {
			return filesystemProviderErr("failed deleting file", err)
		}
	}

	return nil
}

func (fsp FilesystemProvider) GetHashes(subfolder string) ([]string, error) {
	folder := fsp.Config.StorageProvider.FilesystemConfig.Folder + "/" + subfolder

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

func (fsp FilesystemProvider) generateFilepath(subfolder string, hash string, regenerate bool) string {
	var folder string
	if regenerate {
		folder = fsp.Config.StorageProvider.FilesystemConfig.RegenerateFolder
	} else {
		folder = fsp.Config.StorageProvider.FilesystemConfig.Folder
	}

	format := fsp.Config.StorageProvider.FilesystemConfig.Format

	return filepath.Join(folder, subfolder, hash+"."+format)
}

func filesystemProviderErr(message string, err error) error {
	return fmt.Errorf("FileSystemProvider: %s: %v", message, err)
}

func (fsp FilesystemProvider) createFolder(subfolder string, regenerate bool) error {
	var folder string
	if regenerate {
		folder = filepath.Join(fsp.Config.StorageProvider.FilesystemConfig.RegenerateFolder, subfolder)
	} else {
		folder = filepath.Join(fsp.Config.StorageProvider.FilesystemConfig.Folder, subfolder)
	}

	_, err := os.Stat(folder)

	if os.IsNotExist(err) {
		err := os.MkdirAll(folder, 0750)
		if err != nil {
			return err
		}
	}

	return nil
}
