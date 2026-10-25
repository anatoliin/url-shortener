package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/anatoliin/url-shortener/internal/repository"
	"github.com/anatoliin/url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
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
				statusCode:   404,
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
			repo := repository.NewURLMemoryRepository(make(map[string]string))
			urlService := service.NewURLService(repo)
			r := chi.NewRouter()
			r.Method(http.MethodPost, "/", NewURLPostHandler(urlService))
			ts := httptest.NewServer(r)
			defer ts.Close()

			resp, bodyData := testRequest(t, ts, "POST", tt.request, strings.NewReader(tt.body))
			defer resp.Body.Close()
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			if tt.want.checkContent {
				assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
				parsed, err := url.Parse(bodyData)
				require.NoError(t, err)
				assert.Equal(t, "http", parsed.Scheme)

				tsURL, err := url.Parse(ts.URL)
				require.NoError(t, err)
				assert.Equal(t, tsURL.Host, parsed.Host)

				assert.NotEmpty(t, parsed.Path)
			}
		})
	}
}
