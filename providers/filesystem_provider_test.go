package providers

import (
	"reflect"
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

func TestGetConfig(t *testing.T) {
	config := filesystemProvider.GetConfig()

	if !reflect.DeepEqual(testConfig.StorageProvider, config) {
		t.Errorf("expected storage provider config: %v, actual storage provider config: %v", testConfig.StorageProvider, config)
	}
}

func TestGenerateFilepath(t *testing.T) {
	generateFilepath := filesystemProvider.GenerateFilepath("testHash1", false)

	expectedFilepath := "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T/testHash1.json"

	if expectedFilepath != generateFilepath {
		t.Errorf("expected generate filepath: %v, actual generate filepath: %v", expectedFilepath, generateFilepath)
	}

	filesystemProvider.Config.Regenerate = true

	regenerateFilepath := filesystemProvider.GenerateFilepath("testHash2", true)

	expectedFilepath = "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T/regenerate/testHash2.json"

	if expectedFilepath != regenerateFilepath {
		t.Errorf("expected regenerate filepath: %v, actual regenerate filepath: %v", expectedFilepath, regenerateFilepath)
	}
}

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

// reusable functions

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
