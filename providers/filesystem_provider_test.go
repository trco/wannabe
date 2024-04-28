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
	Wannabes: map[string]config.Wannabe{
		"test.api.com": {},
	},
}

var filesystemProvider = FilesystemProvider{
	Config: testConfig,
}

var testRecord = []byte{53}

func TestInsertAndReadRecord(t *testing.T) {
	_ = filesystemProvider.CreateFolders()
	_ = filesystemProvider.InsertRecords([][]byte{testRecord}, []string{"testHash3"}, "test.api.com")

	records, _ := filesystemProvider.ReadRecords([]string{"testHash3"}, "test.api.com")

	if len(records) == 0 {
		t.Errorf("failed inserting and reading records")
	}
}

func TestDeleteRecord(t *testing.T) {
	_ = filesystemProvider.CreateFolders()
	filesystemProvider.InsertRecords([][]byte{testRecord}, []string{"testHash4"}, "test.api.com")

	filesystemProvider.DeleteRecords([]string{"testHash4"}, "test.api.com")

	fileContents, _ := filesystemProvider.ReadRecords([]string{"testHash4"}, "test.api.com")

	if len(fileContents) != 0 {
		t.Errorf("failed deleting record")
	}

	filesystemProvider.DeleteRecords([]string{"testHash4"}, "test.api.com")
}

func TestGetHashes(t *testing.T) {
	_ = filesystemProvider.CreateFolders()
	_ = filesystemProvider.InsertRecords([][]byte{testRecord}, []string{"testHash5"}, "test.api.com")
	_ = filesystemProvider.InsertRecords([][]byte{testRecord}, []string{"testHash6"}, "test.api.com")
	_ = filesystemProvider.InsertRecords([][]byte{testRecord}, []string{"testHash7"}, "test.api.com")

	hashes, _ := filesystemProvider.GetHashes("test.api.com")

	expectedHashes := []string{"testHash5", "testHash6", "testHash7"}

	for _, hash := range expectedHashes {
		if !contains(hashes, hash) {
			t.Errorf("expected hashes: %v does not contain hash: %v", expectedHashes, hash)
		}
	}
}

func TestGenerateFilepath(t *testing.T) {
	generateFilepath := filesystemProvider.generateFilepath("testHash1", "test.api.com")

	expectedFilepath := "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T/test.api.com/testHash1.json"

	if expectedFilepath != generateFilepath {
		t.Errorf("expected generate filepath: %v, actual generate filepath: %v", expectedFilepath, generateFilepath)
	}
}

func TestGenerateFilepathRegenerate(t *testing.T) {
	filesystemProvider.Config.StorageProvider.Regenerate = true

	regenerateFilepath := filesystemProvider.generateFilepathRegenerate("testHash2", "test.api.com")

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
