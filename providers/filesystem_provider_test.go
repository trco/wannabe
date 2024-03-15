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
			RegenerateFolder: "",
			Format:           "json",
		},
	},
}

var filesystemProvider = FilesystemProvider{
	Config: testConfig.StorageProvider,
}

var testRecord = []byte{53}

func TestGenerateFilepath(t *testing.T) {
	filepath := filesystemProvider.GenerateFilepath("testHash1", testConfig.StorageProvider.Regenerate)

	expectedFilepath := "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T/testHash1.json"

	if expectedFilepath != filepath {
		t.Errorf("Expected filepath: %v, Actual filepath: %v", expectedFilepath, filepath)
	}
}

func TestFilesystemProvider(t *testing.T) {
	// Add record
	filesystemProvider.InsertRecords([]string{"testHash2"}, [][]byte{testRecord})

	// ReadRecord
	fileContent, _ := filesystemProvider.ReadRecords([]string{"testHash2"})

	expectedFileContent := []byte{53}

	if !reflect.DeepEqual(expectedFileContent, fileContent) {
		t.Errorf("Expected file content: %v, Actual file content: %v", expectedFileContent, fileContent)
	}
}
