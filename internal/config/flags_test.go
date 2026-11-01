package config

import (
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
