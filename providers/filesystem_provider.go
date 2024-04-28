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

func (fsp FilesystemProvider) CreateFolders() error {
	err := os.Mkdir(fsp.Config.StorageProvider.FilesystemConfig.Folder, 0750)
	if err != nil && !os.IsExist(err) {
		return filesystemProviderErr("failed creating folder defined in 'storageProvider.filesystemConfig.folder' config", err)
	}

	wannabeKeys := fsp.getMapKeys(fsp.Config.Wannabes)

	for _, host := range wannabeKeys {
		err := os.Mkdir(filepath.Join(fsp.Config.StorageProvider.FilesystemConfig.Folder, host), 0750)
		if err != nil && !os.IsExist(err) {
			message := fmt.Sprintf("failed creating folder %v", filepath.Join(fsp.Config.StorageProvider.FilesystemConfig.Folder, host))
			return filesystemProviderErr(message, err)
		}
	}

	if fsp.Config.StorageProvider.Regenerate {
		err = os.Mkdir(fsp.Config.StorageProvider.FilesystemConfig.RegenerateFolder, 0750)
		if err != nil && !os.IsExist(err) {
			return filesystemProviderErr("failed creating regenerate folder or missing 'storageProvider.filesystemConfig.regenerateFolder' config", err)
		}

		for _, host := range wannabeKeys {
			err := os.Mkdir(filepath.Join(fsp.Config.StorageProvider.FilesystemConfig.RegenerateFolder, host), 0750)
			if err != nil && !os.IsExist(err) {
				message := fmt.Sprintf("failed creating folder %v", filepath.Join(fsp.Config.StorageProvider.FilesystemConfig.RegenerateFolder, host))
				return filesystemProviderErr(message, err)
			}
		}
	}

	return nil
}

// TODO ? bulk read using goroutines and channels
func (fsp FilesystemProvider) ReadRecords(hashes []string, subfolder string) ([][]byte, error) {
	var records [][]byte

	// TODO all or nothing ?
	for _, hash := range hashes {
		filepath := fsp.generateFilepath(hash, subfolder)

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
	// TODO all or nothing ?
	for index, record := range records {
		filepath := fsp.generateFilepath(hashes[index], subfolder)
		if fsp.Config.StorageProvider.Regenerate {
			filepath = fsp.generateFilepathRegenerate(hashes[index], subfolder)
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
func (fsp FilesystemProvider) DeleteRecords(hashes []string, subfolder string) error {
	for _, hash := range hashes {
		filepath := fsp.generateFilepath(hash, subfolder)

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

func (fsp FilesystemProvider) getMapKeys(wannabes map[string]config.Wannabe) []string {
	keys := make([]string, 0, len(wannabes))

	for key := range wannabes {
		keys = append(keys, key)
	}

	return keys
}

func (fsp FilesystemProvider) generateFilepath(hash string, subfolder string) string {
	folder := fsp.Config.StorageProvider.FilesystemConfig.Folder
	format := fsp.Config.StorageProvider.FilesystemConfig.Format

	return folder + "/" + subfolder + "/" + hash + "." + format
}

func (fsp FilesystemProvider) generateFilepathRegenerate(hash string, subfolder string) string {
	folder := fsp.Config.StorageProvider.FilesystemConfig.RegenerateFolder
	format := fsp.Config.StorageProvider.FilesystemConfig.Format

	return folder + "/" + subfolder + "/" + hash + "." + format
}

func filesystemProviderErr(message string, err error) error {
	return fmt.Errorf("FileSystemProvider: %s: %v", message, err)
}
