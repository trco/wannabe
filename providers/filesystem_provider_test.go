package providers

import (
	"os"
	"testing"
	"wannabe/types"
)

var testConfig = types.Config{
	StorageProvider: types.StorageProvider{
		Type: "filesystem",
		FilesystemConfig: types.FilesystemConfig{
			Folder:           "/tmp/wannabe",
			RegenerateFolder: "/tmp/wannabe/regenerate",
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
	_ = filesystemProvider.InsertRecords("test.api.com", []string{"testHash3"}, [][]byte{testRecord}, false)

	records, _ := filesystemProvider.ReadRecords("test.api.com", []string{"testHash3"})

	if len(records) == 0 {
		t.Errorf("failed inserting and reading records")
	}
}

func TestDeleteRecord(t *testing.T) {
	filesystemProvider.InsertRecords("test.api.com", []string{"testHash4"}, [][]byte{testRecord}, false)

	filesystemProvider.DeleteRecords("test.api.com", []string{"testHash4"})

	fileContents, _ := filesystemProvider.ReadRecords("test.api.com", []string{"testHash4"})

	if len(fileContents) != 0 {
		t.Errorf("failed deleting record")
	}

	filesystemProvider.DeleteRecords("test.api.com", []string{"testHash4"})
}

func TestGetHashes(t *testing.T) {
	_ = filesystemProvider.InsertRecords("test.api.com", []string{"testHash5"}, [][]byte{testRecord}, false)
	_ = filesystemProvider.InsertRecords("test.api.com", []string{"testHash6"}, [][]byte{testRecord}, false)
	_ = filesystemProvider.InsertRecords("test.api.com", []string{"testHash7"}, [][]byte{testRecord}, false)

	hashes, _ := filesystemProvider.GetHashes("test.api.com")

	expectedHashes := []string{"testHash5", "testHash6", "testHash7"}

	for _, hash := range expectedHashes {
		if !contains(hashes, hash) {
			t.Errorf("expected hashes: %v does not contain hash: %v", expectedHashes, hash)
		}
	}
}

func TestGetHostsAndHashes(t *testing.T) {
	_ = filesystemProvider.InsertRecords("test2.api.com", []string{"testHash8"}, [][]byte{testRecord}, false)
	_ = filesystemProvider.InsertRecords("test2.api.com", []string{"testHash9"}, [][]byte{testRecord}, false)
	_ = filesystemProvider.InsertRecords("test3.api.com", []string{"testHash10"}, [][]byte{testRecord}, false)

	hostsAndHashes, _ := filesystemProvider.GetHostsAndHashes()

	if len(hostsAndHashes) == 0 {
		t.Errorf("failed getting hosts and hashes")
	}
}

func TestGenerateFilepath(t *testing.T) {
	isRegenerate := false
	folder := testConfig.StorageProvider.FilesystemConfig.Folder
	subfolder := "test.api.com"
	hash := "testHash1"

	generateFilepath := filesystemProvider.generateFilepath(subfolder, hash, isRegenerate)

	expectedFilepath := folder + "/" + subfolder + "/" + hash + ".json"

	if expectedFilepath != generateFilepath {
		t.Errorf("expected generate filepath: %v, actual generate filepath: %v", expectedFilepath, generateFilepath)
	}
}

func TestGenerateFilepathRegenerate(t *testing.T) {
	isRegenerate := true
	regenerateFolder := testConfig.StorageProvider.FilesystemConfig.RegenerateFolder
	subfolder := "test.api.com"
	hash := "testHash2"

	regenerateFilepath := filesystemProvider.generateFilepath(subfolder, hash, isRegenerate)

	expectedFilepath := regenerateFolder + "/" + subfolder + "/" + hash + ".json"

	if expectedFilepath != regenerateFilepath {
		t.Errorf("expected regenerate filepath: %v, actual regenerate filepath: %v", expectedFilepath, regenerateFilepath)
	}

	isRegenerate = false
}

func TestCreateFolder(t *testing.T) {
	folder := testConfig.StorageProvider.FilesystemConfig.Folder
	subfolder := "test.subfolder.com"
	isRegenerate := false

	_ = filesystemProvider.createFolder(subfolder, isRegenerate)

	fileInfo, _ := os.Stat(folder + "/" + subfolder)

	if fileInfo.IsDir() == false {
		t.Errorf("failed creating subfolder: %v in folder: %v", subfolder, folder)
	}

	isRegenerate = true
	rengenerateFolder := testConfig.StorageProvider.FilesystemConfig.RegenerateFolder

	_ = filesystemProvider.createFolder(subfolder, isRegenerate)

	fileInfo, _ = os.Stat(rengenerateFolder + "/" + subfolder)

	if fileInfo.IsDir() == false {
		t.Errorf("failed creating subfolder: %v in regenerate folder: %v", subfolder, rengenerateFolder)
	}
}

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func createFolders() {
	_ = os.Mkdir(testConfig.StorageProvider.FilesystemConfig.Folder, 0755)
	_ = os.Mkdir(testConfig.StorageProvider.FilesystemConfig.RegenerateFolder, 0755)
}

func deleteFolders() {
	_ = os.RemoveAll(testConfig.StorageProvider.FilesystemConfig.Folder)
	_ = os.RemoveAll(testConfig.StorageProvider.FilesystemConfig.RegenerateFolder)
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()

	os.Exit(code)
}

func setup() {
	createFolders()
}

func teardown() {
	deleteFolders()
}
