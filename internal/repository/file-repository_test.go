package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewURLFileRepository(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		id          string
		wantURL     string
		wantErr     error
	}{
		{
			name:        "loads URL",
			fileContent: `[{"uuid":"1","short_url":"abc123","original_url":"https://example.com"}]`,
			id:          "abc123",
			wantURL:     "https://example.com",
		},
		{
			name:    "empty file",
			id:      "missing",
			wantErr: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := filepath.Join(t.TempDir(), "urls.json")
			if tt.fileContent != "" {
				require.NoError(t, os.WriteFile(filename, []byte(tt.fileContent), 0600))
			}

			repo, err := NewURLFileRepository(filename)
			require.NoError(t, err)

			got, err := repo.FindByID(tt.id)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantURL, got)
		})
	}
}

func TestURLFileRepository_CloseWritesCurrentData(t *testing.T) {
	tests := []struct {
		name         string
		initialData  string
		id           string
		originalURL  string
		wantURLs     map[string]string
		wantRecordID []string
	}{
		{
			name:        "writes new URL to empty file",
			id:          "new",
			originalURL: "https://new.example",
			wantURLs:    map[string]string{"new": "https://new.example"},
			wantRecordID: []string{
				"1",
			},
		},
		{
			name:        "preserves and appends URL",
			initialData: `[{"uuid":"1","short_url":"old","original_url":"https://old.example"}]`,
			id:          "new",
			originalURL: "https://new.example",
			wantURLs: map[string]string{
				"old": "https://old.example",
				"new": "https://new.example",
			},
			wantRecordID: []string{"1", "2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := filepath.Join(t.TempDir(), "urls.json")
			if tt.initialData != "" {
				require.NoError(t, os.WriteFile(filename, []byte(tt.initialData), 0600))
			}

			repo, err := NewURLFileRepository(filename)
			require.NoError(t, err)
			require.NoError(t, repo.Save(tt.id, tt.originalURL))
			require.NoError(t, repo.Close())

			data, err := os.ReadFile(filename)
			require.NoError(t, err)
			var records []struct {
				ID string `json:"uuid"`
			}
			require.NoError(t, json.Unmarshal(data, &records))
			require.Len(t, records, len(tt.wantRecordID))
			recordIDs := make([]string, len(records))
			for i, record := range records {
				recordIDs[i] = record.ID
			}
			require.ElementsMatch(t, tt.wantRecordID, recordIDs)

			reopened, err := NewURLFileRepository(filename)
			require.NoError(t, err)
			for id, wantURL := range tt.wantURLs {
				got, err := reopened.FindByID(id)
				require.NoError(t, err)
				require.Equal(t, wantURL, got)
			}
		})
	}
}
