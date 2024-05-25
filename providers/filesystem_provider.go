package providers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"wannabe/types"
)

type FilesystemProvider struct {
	Config types.Config
}

func (fsp FilesystemProvider) ReadRecords(subfolder string, hashes []string) ([][]byte, error) {
	var records [][]byte

	// TODO all or nothing ?
	for _, hash := range hashes {
		filepath := fsp.generateFilepath(subfolder, hash, false)

		_, err := os.Stat(filepath)
		if os.IsNotExist(err) {
			continue
		} else if err != nil {
			return nil, filesystemProviderErr("failed checking if file exists "+filepath, err)
		}

		file, err := os.Open(filepath)
		if err != nil {
			return nil, filesystemProviderErr("failed opening file "+filepath, err)
		}
		defer file.Close()

		record, err := io.ReadAll(file)
		if err != nil {
			return nil, filesystemProviderErr("failed reading file "+filepath, err)
		}

		records = append(records, record)
	}

	return records, nil
}

func (fsp FilesystemProvider) InsertRecords(subfolder string, hashes []string, records [][]byte) error {
	isRegenerate := fsp.Config.StorageProvider.Regenerate

	err := fsp.createFolder(subfolder, isRegenerate)
	if err != nil {
		return filesystemProviderErr("failed creating folder "+subfolder, err)
	}

	// TODO all or nothing ?
	for index, record := range records {
		filepath := fsp.generateFilepath(subfolder, hashes[index], isRegenerate)

		_, err := os.Create(filepath)
		if err != nil {
			return filesystemProviderErr("failed creating file "+filepath, err)
		}

		err = os.WriteFile(filepath, record, 0644)
		if err != nil {
			return filesystemProviderErr("failed writing file "+filepath, err)
		}
	}

	return nil
}

func (fsp FilesystemProvider) DeleteRecords(subfolder string, hashes []string) error {
	for _, hash := range hashes {
		filepath := fsp.generateFilepath(subfolder, hash, false)

		err := os.Remove(filepath)
		if err != nil {
			return filesystemProviderErr("failed deleting file "+filepath, err)
		}
	}

	return nil
}

func (fsp FilesystemProvider) GetHashes(subfolder string) ([]string, error) {
	folder := fsp.Config.StorageProvider.FilesystemConfig.Folder + "/" + subfolder

	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, filesystemProviderErr("failed reading folder "+fsp.Config.StorageProvider.FilesystemConfig.Folder, err)
	}

	hashes := fsp.getHashesFromFiles(files)

	return hashes, nil
}

func (fsp FilesystemProvider) GetHostsAndHashes() ([]HostAndHashes, error) {
	var hostHashes []HostAndHashes

	folder := fsp.Config.StorageProvider.FilesystemConfig.Folder

	subfolders, err := os.ReadDir(folder)
	if err != nil {
		return nil, filesystemProviderErr("failed reading folder "+folder, err)
	}

	for _, subfolder := range subfolders {
		files, err := os.ReadDir(folder + "/" + subfolder.Name())
		if err != nil {
			return nil, filesystemProviderErr("failed reading subfolder "+subfolder.Name(), err)
		}

		hashes := fsp.getHashesFromFiles(files)

		hostHashes = append(hostHashes, HostAndHashes{
			Host:      subfolder.Name(),
			HashCount: len(hashes),
			Hashes:    hashes,
		})

	}

	return hostHashes, nil
}

func (fsp FilesystemProvider) generateFilepath(subfolder string, hash string, isRegenerate bool) string {
	var folder string
	if isRegenerate {
		folder = fsp.Config.StorageProvider.FilesystemConfig.RegenerateFolder
	} else {
		folder = fsp.Config.StorageProvider.FilesystemConfig.Folder
	}

	format := fsp.Config.StorageProvider.FilesystemConfig.Format

	return filepath.Join(folder, subfolder, hash+"."+format)
}

func (fsp FilesystemProvider) createFolder(subfolder string, isRegenerate bool) error {
	var folder string
	if isRegenerate {
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

func (fsp FilesystemProvider) getHashesFromFiles(files []os.DirEntry) []string {
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

	return hashes
}

func filesystemProviderErr(message string, err error) error {
	return fmt.Errorf("FileSystemProvider: %s: %v", message, err)
}
