package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBaseURL_Set(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    string
		wantErr bool
	}{
		{
			name:  "https with path",
			value: "https://example.com/short",
			want:  "https://example.com/short",
		},
		{
			name:  "http",
			value: "http://localhost:8080",
			want:  "http://localhost:8080",
		},
		{
			name:    "empty",
			value:   "",
			wantErr: true,
		},
		{
			name:    "no scheme",
			value:   "localhost:8080",
			wantErr: true,
		},
		{
			name:    "ftp scheme",
			value:   "ftp://example.com",
			wantErr: true,
		},
		{
			name:    "scheme only",
			value:   "http://",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b baseURL
			err := b.Set(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, b.String())
			assert.Equal(t, tt.want, b.URL().String())
		})
	}
}

func TestBaseURL_DefaultConstant(t *testing.T) {
	var b baseURL
	require.NoError(t, b.Set(defaultBaseURL))
	assert.Equal(t, defaultBaseURL, b.String())
}

func TestInitOptions(t *testing.T) {
	tests := []struct {
		name                string
		args                []string
		serverAddr          string
		fileStoragePath     string
		baseURL             string
		wantAddress         string
		wantFileStoragePath string
		wantBaseURL         string
	}{
		{
			name:                "defaults",
			wantAddress:         ":8080",
			wantFileStoragePath: "data.json",
			wantBaseURL:         defaultBaseURL,
		},
		{
			name:                "command line flags",
			args:                []string{"-a", ":9090", "-f", "urls.json", "-b", "https://example.com/short"},
			wantAddress:         ":9090",
			wantFileStoragePath: "urls.json",
			wantBaseURL:         "https://example.com/short",
		},
		{
			name:                "environment variables",
			serverAddr:          ":7070",
			fileStoragePath:     "env-urls.json",
			baseURL:             "https://env.example.com",
			wantAddress:         ":7070",
			wantFileStoragePath: "env-urls.json",
			wantBaseURL:         "https://env.example.com",
		},
		{
			name:                "server address environment variable only",
			serverAddr:          ":7070",
			wantAddress:         ":7070",
			wantFileStoragePath: "data.json",
			wantBaseURL:         defaultBaseURL,
		},
		{
			name:                "base URL environment variable only",
			baseURL:             "https://env.example.com",
			wantAddress:         ":8080",
			wantFileStoragePath: "data.json",
			wantBaseURL:         "https://env.example.com",
		},
		{
			name:                "file storage path environment variable only",
			fileStoragePath:     "env-only.json",
			wantAddress:         ":8080",
			wantFileStoragePath: "env-only.json",
			wantBaseURL:         defaultBaseURL,
		},
		{
			name:                "environment variables override flags",
			args:                []string{"-a", ":9090", "-f", "flag-urls.json", "-b", "https://flag.example.com"},
			serverAddr:          ":7070",
			fileStoragePath:     "env-urls.json",
			baseURL:             "https://env.example.com",
			wantAddress:         ":7070",
			wantFileStoragePath: "env-urls.json",
			wantBaseURL:         "https://env.example.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldCommandLine, oldArgs := flag.CommandLine, os.Args
			flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
			os.Args = append([]string{"test"}, tt.args...)
			t.Cleanup(func() {
				flag.CommandLine = oldCommandLine
				os.Args = oldArgs
			})

			t.Setenv("SERVER_ADDRESS", tt.serverAddr)
			t.Setenv("FILE_STORAGE_PATH", tt.fileStoragePath)
			t.Setenv("BASE_URL", tt.baseURL)

			got := InitOptions()
			assert.Equal(t, tt.wantAddress, got.RunAddress)
			assert.Equal(t, tt.wantFileStoragePath, got.FileStoragePath)
			assert.Equal(t, tt.wantBaseURL, got.RedirectPrefix.String())
		})
	}
}
