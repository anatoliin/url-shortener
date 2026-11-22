package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestURLMemoryRepository(t *testing.T) {
	tests := []struct {
		name        string
		initialData map[string]string
		id          string
		url         string
		save        bool
		wantURL     string
		wantSaveErr error
		wantFindErr error
	}{
		{
			name:        "save and find URL",
			initialData: map[string]string{},
			id:          "abc123",
			url:         "https://example.com",
			save:        true,
			wantURL:     "https://example.com",
		},
		{
			name:        "find existing URL",
			initialData: map[string]string{"abc123": "https://example.com"},
			id:          "abc123",
			wantURL:     "https://example.com",
		},
		{
			name:        "missing URL",
			initialData: map[string]string{},
			id:          "missing",
			wantFindErr: ErrNotFound,
		},
		{
			name:        "duplicate ID",
			initialData: map[string]string{"abc123": "https://example.com"},
			id:          "abc123",
			url:         "https://other.example",
			save:        true,
			wantURL:     "https://example.com",
			wantSaveErr: ErrCollision,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewURLMemoryRepository(tt.initialData)

			if tt.save {
				err := repo.Save(tt.id, tt.url)
				if tt.wantSaveErr != nil {
					require.ErrorIs(t, err, tt.wantSaveErr)
				} else {
					require.NoError(t, err)
				}
			}

			got, err := repo.FindByID(tt.id)
			if tt.wantFindErr != nil {
				require.ErrorIs(t, err, tt.wantFindErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantURL, got)
		})
	}
}
