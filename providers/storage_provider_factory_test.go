package providers

import (
	"reflect"
	"testing"
	"wannabe/types"
)

func TestStorageProviderFactory(t *testing.T) {
	tests := []struct {
		name   string
		config types.Config
		want   interface{}
		err    bool
	}{
		{
			name: "invalid filesystem config",
			config: types.Config{
				StorageProvider: types.StorageProvider{
					Type: "test",
				},
			},
			want: nil,
			err:  true,
		},
		{
			name: "valid filesystem config",
			config: types.Config{
				StorageProvider: types.StorageProvider{
					Type: "filesystem",
					FilesystemConfig: types.FilesystemConfig{
						Folder:           "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T",
						RegenerateFolder: "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T/regenerate",
						Format:           "json",
					},
				},
			},
			want: FilesystemProvider{
				Config: types.Config{
					StorageProvider: types.StorageProvider{
						Type: "filesystem",
						FilesystemConfig: types.FilesystemConfig{
							Folder:           "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T",
							RegenerateFolder: "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T/regenerate",
							Format:           "json",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StorageProviderFactory(tt.config)

			if (err != nil) != tt.err {
				t.Errorf("StorageProviderFactory() error = %v, wantErr %v", err, tt.err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StorageProviderFactory() = %v, want %v", got, tt.want)
			}
		})
	}
}
