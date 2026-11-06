package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/anatoliin/url-shortener/internal/repository"
	"github.com/anatoliin/url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestURLShortenHandler_ServeHTTP(t *testing.T) {
	type want struct {
		statusCode   int
		contentType  string
		checkContent bool
	}
	tests := []struct {
		name string
		path string
		body string
		want want
	}{
		{
			name: "correct https url",
			path: "/api/shorten",
			body: `{"url":"https://google.com"}`,
			want: want{
				statusCode:   http.StatusCreated,
				contentType:  "text/plain",
				checkContent: true,
			},
		},
		{
			name: "correct http url",
			path: "/api/shorten",
			body: `{"url":"http://google.com"}`,
			want: want{
				statusCode:   http.StatusCreated,
				contentType:  "text/plain",
				checkContent: true,
			},
		},
		{
			name: "incorrect url",
			path: "/api/shorten",
			body: `{"url":"google.com"}`,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "unsupported scheme",
			path: "/api/shorten",
			body: `{"url":"ftp://google.com"}`,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "wrong path",
			path: "/foo",
			body: `{"url":"http://google.com"}`,
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "empty body",
			path: "/api/shorten",
			body: "",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "malformed JSON",
			path: "/api/shorten",
			body: `{"url":`,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "empty URL",
			path: "/api/shorten",
			body: `{"url":""}`,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "correct URL with path and query",
			path: "/api/shorten",
			body: `{"url":"https://google.com/search?q=1"}`,
			want: want{
				statusCode:   http.StatusCreated,
				contentType:  "text/plain",
				checkContent: true,
			},
		},
		{
			name: "oversized body",
			path: "/api/shorten",
			body: `{"url":"https://example.com/` + strings.Repeat("a", maxBodySize) + `"}`,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}

	redirectURLStr := "https://test.com/short-url"
	redirectURL, err := url.Parse(redirectURLStr)
	assert.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repository.NewURLMemoryRepository(make(map[string]string))
			urlService := service.NewURLService(repo)
			r := chi.NewRouter()
			r.Method(http.MethodPost, "/api/shorten", NewURLShortenHandle(urlService, redirectURL))
			ts := httptest.NewServer(r)
			defer ts.Close()

			resp, bodyData := testRequest(t, ts, http.MethodPost, tt.path, strings.NewReader(tt.body))
			defer resp.Body.Close()
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			if tt.want.checkContent {
				assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
				var response ShortenResponse
				assert.NoError(t, json.Unmarshal([]byte(bodyData), &response))
				assert.True(t, strings.HasPrefix(response.Result, redirectURLStr+"/"))
			}
		})
	}
}
