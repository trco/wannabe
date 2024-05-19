package providers

import (
	"testing"
	"wannabe/types"
)

var testConfig = types.Config{
	StorageProvider: types.StorageProvider{
		Type:       "filesystem",
		Regenerate: false,
		FilesystemConfig: types.FilesystemConfig{
			Folder:           "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T",
			RegenerateFolder: "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T/regenerate",
			Format:           "json",
		},
	},
	Wannabes: map[string]types.Wannabe{
		"test.api.com": {},
	},
}

var filesystemProvider = FilesystemProvider{
	Config: testConfig,
}

var testRecord = []byte{53}

func TestInsertAndReadRecord(t *testing.T) {
	_ = filesystemProvider.InsertRecords("test.api.com", []string{"testHash3"}, [][]byte{testRecord})

	records, _ := filesystemProvider.ReadRecords("test.api.com", []string{"testHash3"})

	if len(records) == 0 {
		t.Errorf("failed inserting and reading records")
	}
}

func TestDeleteRecord(t *testing.T) {
	filesystemProvider.InsertRecords("test.api.com", []string{"testHash4"}, [][]byte{testRecord})

	filesystemProvider.DeleteRecords("test.api.com", []string{"testHash4"})

	fileContents, _ := filesystemProvider.ReadRecords("test.api.com", []string{"testHash4"})

	if len(fileContents) != 0 {
		t.Errorf("failed deleting record")
	}

	filesystemProvider.DeleteRecords("test.api.com", []string{"testHash4"})
}

func TestGetHashes(t *testing.T) {
	_ = filesystemProvider.InsertRecords("test.api.com", []string{"testHash5"}, [][]byte{testRecord})
	_ = filesystemProvider.InsertRecords("test.api.com", []string{"testHash6"}, [][]byte{testRecord})
	_ = filesystemProvider.InsertRecords("test.api.com", []string{"testHash7"}, [][]byte{testRecord})

	hashes, _ := filesystemProvider.GetHashes("test.api.com")

	expectedHashes := []string{"testHash5", "testHash6", "testHash7"}

	for _, hash := range expectedHashes {
		if !contains(hashes, hash) {
			t.Errorf("expected hashes: %v does not contain hash: %v", expectedHashes, hash)
		}
	}
}

func TestGenerateFilepath(t *testing.T) {
	regenerate := filesystemProvider.Config.StorageProvider.Regenerate

	generateFilepath := filesystemProvider.generateFilepath("test.api.com", "testHash1", regenerate)

	expectedFilepath := "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T/test.api.com/testHash1.json"

	if expectedFilepath != generateFilepath {
		t.Errorf("expected generate filepath: %v, actual generate filepath: %v", expectedFilepath, generateFilepath)
	}
}

func TestGenerateFilepathRegenerate(t *testing.T) {
	filesystemProvider.Config.StorageProvider.Regenerate = true

	isRegenerate := filesystemProvider.Config.StorageProvider.Regenerate

	regenerateFilepath := filesystemProvider.generateFilepath("test.api.com", "testHash2", isRegenerate)

	expectedFilepath := "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T/regenerate/test.api.com/testHash2.json"

	if expectedFilepath != regenerateFilepath {
		t.Errorf("expected regenerate filepath: %v, actual regenerate filepath: %v", expectedFilepath, regenerateFilepath)
	}

	filesystemProvider.Config.StorageProvider.Regenerate = false
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
