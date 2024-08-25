package storage

import (
	"reflect"
	"testing"

	"github.com/trco/wannabe/internal/config"
)

func TestProviderFactory(t *testing.T) {
	tests := []struct {
		name   string
		config config.StorageProvider
		want   interface{}
		err    bool
	}{
		{
			name: "invalid filesystem config",
			config: config.StorageProvider{
				Type: "test",
			},
			want: nil,
			err:  true,
		},
		{
			name: "valid filesystem config",
			config: config.StorageProvider{
				Type: "filesystem",
				FilesystemConfig: config.FilesystemConfig{
					Folder:           "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T",
					RegenerateFolder: "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T/regenerate",
					Format:           "json",
				},
			},
			want: FilesystemProvider{
				Config: config.StorageProvider{
					Type: "filesystem",
					FilesystemConfig: config.FilesystemConfig{
						Folder:           "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T",
						RegenerateFolder: "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T/regenerate",
						Format:           "json",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProviderFactory(tt.config)

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
