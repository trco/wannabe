package providers

import (
	"testing"
	"wannabe/config"
)

var testConfig = config.Config{
	StorageProvider: config.StorageProvider{
		Type:       "filesystem",
		Regenerate: false,
		FilesystemConfig: config.FilesystemConfig{
			Folder:           "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T",
			RegenerateFolder: "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T/regenerate",
			Format:           "json",
		},
	},
}

var filesystemProvider = FilesystemProvider{
	Config: testConfig.StorageProvider,
}

var testRecord = []byte{53}

func TestInsertAndReadRecord(t *testing.T) {
	_ = filesystemProvider.InsertRecords([]string{"testHash3"}, [][]byte{testRecord})

	records, _ := filesystemProvider.ReadRecords([]string{"testHash3"})

	if len(records) == 0 {
		t.Errorf("failed inserting and reading records")
	}
}

func TestDeleteRecord(t *testing.T) {
	filesystemProvider.InsertRecords([]string{"testHash4"}, [][]byte{testRecord})

	filesystemProvider.DeleteRecords([]string{"testHash4"})

	fileContents, _ := filesystemProvider.ReadRecords([]string{"testHash4"})

	if len(fileContents) != 0 {
		t.Errorf("failed deleting record")
	}

	filesystemProvider.DeleteRecords([]string{"testHash4"})
}

func TestGetHashes(t *testing.T) {
	_ = filesystemProvider.InsertRecords([]string{"testHash5"}, [][]byte{testRecord})
	_ = filesystemProvider.InsertRecords([]string{"testHash6"}, [][]byte{testRecord})
	_ = filesystemProvider.InsertRecords([]string{"testHash7"}, [][]byte{testRecord})

	hashes, _ := filesystemProvider.GetHashes()

	expectedHashes := []string{"testHash5", "testHash6", "testHash7"}

	for _, hash := range expectedHashes {
		if !contains(hashes, hash) {
			t.Errorf("expected hashes: %v does not contain hash: %v", expectedHashes, hash)
		}
	}
}

func TestGenerateFilepath(t *testing.T) {
	generateFilepath := filesystemProvider.generateFilepath("testHash1")

	expectedFilepath := "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T/testHash1.json"

	if expectedFilepath != generateFilepath {
		t.Errorf("expected generate filepath: %v, actual generate filepath: %v", expectedFilepath, generateFilepath)
	}

	filesystemProvider.Config.Regenerate = true

	regenerateFilepath := filesystemProvider.generateFilepath("testHash2")

	expectedFilepath = "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T/regenerate/testHash2.json"

	if expectedFilepath != regenerateFilepath {
		t.Errorf("expected regenerate filepath: %v, actual regenerate filepath: %v", expectedFilepath, regenerateFilepath)
	}

	filesystemProvider.Config.Regenerate = false
}

func TestGetFolder(t *testing.T) {
	filesystemProvider.Config.Regenerate = false

	folder := filesystemProvider.getFolder()
	expectedFolder := "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T"

	if expectedFolder != folder {
		t.Errorf("expected folder: %v, actual folder: %v", expectedFolder, folder)
	}

	filesystemProvider.Config.Regenerate = true

	regenerateFolder := filesystemProvider.getFolder()
	expectedRegenerateFolder := "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T/regenerate"

	if expectedRegenerateFolder != regenerateFolder {
		t.Errorf("expected regenerate folder: %v, actual regenerate folder: %v", expectedRegenerateFolder, regenerateFolder)
	}
}

// reusable functions

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
