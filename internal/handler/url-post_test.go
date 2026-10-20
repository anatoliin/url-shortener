package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/anatoliin/url-shortener/internal/repository"
	"github.com/anatoliin/url-shortener/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURLPostHandler_ServeHTTP(t *testing.T) {
	type want struct {
		statusCode   int
		contentType  string
		checkContent bool
	}
	tests := []struct {
		name    string
		request string
		body    string
		want    want
	}{
		{
			name:    "correct https url test #1",
			request: "/",
			body:    "https://google.com",
			want: want{
				statusCode:   201,
				contentType:  "text/plain",
				checkContent: true,
			},
		},
		{
			name:    "correct http url test #2",
			request: "/",
			body:    "http://google.com",
			want: want{
				statusCode:   201,
				contentType:  "text/plain",
				checkContent: true,
			},
		},
		{
			name:    "incorrect url test #3",
			request: "/",
			body:    "google.com",
			want: want{
				statusCode:   400,
				contentType:  "",
				checkContent: false,
			},
		},
		{
			name:    "incorrect url test #4",
			request: "/",
			body:    "ftp://google.com",
			want: want{
				statusCode:   400,
				contentType:  "",
				checkContent: false,
			},
		},
		{
			name:    "wrong path test #5",
			request: "/foo",
			body:    "http://google.com",
			want: want{
				statusCode:   400,
				contentType:  "",
				checkContent: false,
			},
		},
		{
			name:    "empty body test #6",
			request: "/",
			body:    "",
			want: want{
				statusCode:   400,
				contentType:  "",
				checkContent: false,
			},
		},
		{
			name:    "scheme only test #6",
			request: "/",
			body:    "http://",
			want: want{
				statusCode:   400,
				contentType:  "",
				checkContent: false,
			},
		},
		{
			name:    "correct url with path/query test #7",
			request: "/",
			body:    "https://google.com/search?q=1",
			want: want{
				statusCode:   201,
				contentType:  "text/plain",
				checkContent: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, strings.NewReader(tt.body))
			request.Host = "localhost:8080"
			w := httptest.NewRecorder()
			repo := repository.NewURLMemoryRepository(make(map[string]string))
			urlService := service.NewURLService(repo)
			h := http.Handler(NewURLPostHandler(urlService))
			h.ServeHTTP(w, request)

			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			bodyData, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			if tt.want.checkContent {
				assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
				parsed, err := url.Parse(string(bodyData))
				require.NoError(t, err)
				assert.Equal(t, "http", parsed.Scheme)
				assert.Equal(t, request.Host, parsed.Host)
				assert.NotEmpty(t, parsed.Path)
			}
		})
	}
}
