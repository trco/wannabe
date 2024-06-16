package providers

import (
	"os"
	"reflect"
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
	t.Run("insert and read record", func(t *testing.T) {
		_ = filesystemProvider.InsertRecords("test.api.com", []string{"testHash3"}, [][]byte{testRecord}, false)

		records, _ := filesystemProvider.ReadRecords("test.api.com", []string{"testHash3"})

		want := 1
		got := len(records)

		if got != want {
			t.Errorf("InsertRecords() + ReadRecords() = %v, want %v", got, want)
		}
	})
}

func TestDeleteRecord(t *testing.T) {
	t.Run("delete record", func(t *testing.T) {
		filesystemProvider.InsertRecords("test.api.com", []string{"testHash4"}, [][]byte{testRecord}, false)
		filesystemProvider.DeleteRecords("test.api.com", []string{"testHash4"})
		fileContents, _ := filesystemProvider.ReadRecords("test.api.com", []string{"testHash4"})

		want := 0
		got := len(fileContents)

		if got != want {
			t.Errorf("InsertRecords() + DeleteRecords() results in %v records stored in filesystem, want %v records stored in filesystem", got, want)
		}

		filesystemProvider.DeleteRecords("test.api.com", []string{"testHash4"})
	})
}

func TestGetHashes(t *testing.T) {
	t.Run("get hashes", func(t *testing.T) {
		_ = filesystemProvider.InsertRecords("test.api.com", []string{"testHash5"}, [][]byte{testRecord}, false)
		_ = filesystemProvider.InsertRecords("test.api.com", []string{"testHash6"}, [][]byte{testRecord}, false)
		_ = filesystemProvider.InsertRecords("test.api.com", []string{"testHash7"}, [][]byte{testRecord}, false)

		want := []string{"testHash5", "testHash6", "testHash7"}
		got, _ := filesystemProvider.GetHashes("test.api.com")

		if !reflect.DeepEqual(got, want) {
			t.Errorf("GetHashes() = %v, want %v", got, want)
		}
	})
}

func TestGetHostsAndHashes(t *testing.T) {
	t.Run("get hosts and hashes", func(t *testing.T) {
		_ = filesystemProvider.InsertRecords("test2.api.com", []string{"testHash8"}, [][]byte{testRecord}, false)
		_ = filesystemProvider.InsertRecords("test2.api.com", []string{"testHash9"}, [][]byte{testRecord}, false)
		_ = filesystemProvider.InsertRecords("test3.api.com", []string{"testHash10"}, [][]byte{testRecord}, false)

		hostsAndHashes, _ := filesystemProvider.GetHostsAndHashes()

		want := 3
		got := len(hostsAndHashes)

		if got != want {
			t.Errorf("GetHostsAndHashes() returns %v entries, want %v entries", got, want)
		}
	})
}

func TestGenerateFilepath(t *testing.T) {
	t.Run("generate filepath", func(t *testing.T) {
		isRegenerate := false
		folder := testConfig.StorageProvider.FilesystemConfig.Folder
		subfolder := "test.api.com"
		hash := "testHash1"

		want := folder + "/" + subfolder + "/" + hash + ".json"

		got := filesystemProvider.generateFilepath(subfolder, hash, isRegenerate)

		if got != want {
			t.Errorf("generateFilepath() = %v, want %v", got, want)
		}
	})
}

func TestGenerateFilepathRegenerate(t *testing.T) {
	t.Run("generate filepath regenerate", func(t *testing.T) {
		isRegenerate := true
		regenerateFolder := testConfig.StorageProvider.FilesystemConfig.RegenerateFolder
		subfolder := "test.api.com"
		hash := "testHash2"

		want := regenerateFolder + "/" + subfolder + "/" + hash + ".json"

		got := filesystemProvider.generateFilepath(subfolder, hash, isRegenerate)

		if got != want {
			t.Errorf("generateFilepath() = %v, want %v", got, want)
		}

		isRegenerate = false
	})
}

func TestCreateFolder(t *testing.T) {
	t.Run("create folder", func(t *testing.T) {
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

		want := true
		got := fileInfo.IsDir()

		if got != want {
			t.Errorf("createFolder() = %v, want %v", got, want)
		}
	})

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
