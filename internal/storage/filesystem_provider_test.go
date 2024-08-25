package storage

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/trco/wannabe/internal/config"
)

var testConfig = config.StorageProvider{
	Type: "filesystem",
	FilesystemConfig: config.FilesystemConfig{
		Folder:           "/tmp/wannabe",
		RegenerateFolder: "/tmp/wannabe/regenerate",
		Format:           "json",
	},
}

var testFilesystemProvider = FilesystemProvider{
	Config: testConfig,
}

var testRecord = []byte{53}

func TestInsertAndReadRecord(t *testing.T) {
	t.Run("insert and read record", func(t *testing.T) {
		t.Cleanup(teardown)

		_ = testFilesystemProvider.InsertRecords("test.api.com", []string{"testHash3"}, [][]byte{testRecord}, false)

		records, _ := testFilesystemProvider.ReadRecords("test.api.com", []string{"testHash3"})

		want := 1
		got := len(records)

		if got != want {
			t.Errorf("InsertRecords() + ReadRecords() = %v, want %v", got, want)
		}
	})
}

func TestDeleteRecord(t *testing.T) {
	t.Run("delete record", func(t *testing.T) {
		t.Cleanup(teardown)

		testFilesystemProvider.InsertRecords("test.api.com", []string{"testHash4"}, [][]byte{testRecord}, false)
		testFilesystemProvider.DeleteRecords("test.api.com", []string{"testHash4"})
		fileContents, _ := testFilesystemProvider.ReadRecords("test.api.com", []string{"testHash4"})

		want := 0
		got := len(fileContents)

		if got != want {
			t.Errorf("InsertRecords() + DeleteRecords() results in %v records stored in filesystem, want %v records stored in filesystem", got, want)
		}
	})
}

func TestGetHashes(t *testing.T) {
	t.Run("get hashes", func(t *testing.T) {
		t.Cleanup(teardown)

		_ = testFilesystemProvider.InsertRecords("test.api.com", []string{"testHash5"}, [][]byte{testRecord}, false)
		_ = testFilesystemProvider.InsertRecords("test.api.com", []string{"testHash6"}, [][]byte{testRecord}, false)
		_ = testFilesystemProvider.InsertRecords("test.api.com", []string{"testHash7"}, [][]byte{testRecord}, false)

		want := []string{"testHash5", "testHash6", "testHash7"}
		got, _ := testFilesystemProvider.GetHashes("test.api.com")

		if !reflect.DeepEqual(got, want) {
			t.Errorf("GetHashes() = %v, want %v", got, want)
		}
	})
}

func TestGetHostsAndHashes(t *testing.T) {
	t.Run("get hosts and hashes", func(t *testing.T) {
		t.Cleanup(teardown)

		_ = testFilesystemProvider.InsertRecords("test2.api.com", []string{"testHash8"}, [][]byte{testRecord}, false)
		_ = testFilesystemProvider.InsertRecords("test2.api.com", []string{"testHash9"}, [][]byte{testRecord}, false)
		_ = testFilesystemProvider.InsertRecords("test3.api.com", []string{"testHash10"}, [][]byte{testRecord}, false)

		hostsAndHashes, _ := testFilesystemProvider.GetHostsAndHashes()

		fmt.Println(hostsAndHashes)

		want := 2
		got := len(hostsAndHashes)

		if got != want {
			t.Errorf("GetHostsAndHashes() returns %v entries, want %v entries", got, want)
		}
	})
}

func TestGenerateFilepath(t *testing.T) {
	t.Run("generate filepath", func(t *testing.T) {
		isRegenerate := false
		folder := testConfig.FilesystemConfig.Folder
		subfolder := "test.api.com"
		hash := "testHash1"

		want := folder + "/" + subfolder + "/" + hash + ".json"

		got := testFilesystemProvider.generateFilepath(subfolder, hash, isRegenerate)

		if got != want {
			t.Errorf("generateFilepath() = %v, want %v", got, want)
		}
	})
}

func TestGenerateFilepathRegenerate(t *testing.T) {
	t.Run("generate filepath regenerate", func(t *testing.T) {
		isRegenerate := true
		regenerateFolder := testConfig.FilesystemConfig.RegenerateFolder
		subfolder := "test.api.com"
		hash := "testHash2"

		want := regenerateFolder + "/" + subfolder + "/" + hash + ".json"

		got := testFilesystemProvider.generateFilepath(subfolder, hash, isRegenerate)

		if got != want {
			t.Errorf("generateFilepath() = %v, want %v", got, want)
		}

		isRegenerate = false
	})
}

func TestCreateFolder(t *testing.T) {
	t.Run("create folder", func(t *testing.T) {
		t.Cleanup(teardown)

		folder := testConfig.FilesystemConfig.Folder
		subfolder := "test.subfolder.com"
		isRegenerate := false

		_ = testFilesystemProvider.createFolder(subfolder, isRegenerate)

		fileInfo, _ := os.Stat(folder + "/" + subfolder)

		if fileInfo.IsDir() == false {
			t.Errorf("failed creating subfolder: %v in folder: %v", subfolder, folder)
		}

		isRegenerate = true
		rengenerateFolder := testConfig.FilesystemConfig.RegenerateFolder

		_ = testFilesystemProvider.createFolder(subfolder, isRegenerate)

		fileInfo, _ = os.Stat(rengenerateFolder + "/" + subfolder)

		want := true
		got := fileInfo.IsDir()

		if got != want {
			t.Errorf("createFolder() = %v, want %v", got, want)
		}
	})

}

func createFolders() {
	_ = os.Mkdir(testConfig.FilesystemConfig.Folder, 0755)
	_ = os.Mkdir(testConfig.FilesystemConfig.RegenerateFolder, 0755)
}

func deleteFolders() {
	_ = os.RemoveAll(testConfig.FilesystemConfig.Folder)
	_ = os.RemoveAll(testConfig.FilesystemConfig.RegenerateFolder)
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
