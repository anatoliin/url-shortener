package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/anatoliin/url-shortener/internal/repository"
	"github.com/anatoliin/url-shortener/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestURLGetHandler_ServeHTTP(t *testing.T) {
	type want struct {
		statusCode int
		location   string
	}
	tests := []struct {
		name    string
		request string
		storage map[string]string
		want    want
	}{
		{
			name:    "existing https url test #1",
			request: "/aa",
			storage: map[string]string{"aa": "https://gobyexample.com/maps", "bb": "https://gobyexample.com/maps2"},
			want: want{
				statusCode: 307,
				location:   "https://gobyexample.com/maps",
			},
		},
		{
			name:    "existing http url test #1.2",
			request: "/aa",
			storage: map[string]string{"aa": "http://gobyexample.com/maps", "bb": "https://gobyexample.com/maps2"},
			want: want{
				statusCode: 307,
				location:   "http://gobyexample.com/maps",
			},
		},
		{
			name:    "existing url with query test #2",
			request: "/bb",
			storage: map[string]string{"aa": "https://gobyexample.com/maps", "bb": "https://gobyexample.com/maps?b=1&c=2"},
			want: want{
				statusCode: 307,
				location:   "https://gobyexample.com/maps?b=1&c=2",
			},
		},
		{
			name:    "absent url test #3",
			request: "/aA",
			storage: map[string]string{"aa": "https://gobyexample.com/maps", "bb": "https://gobyexample.com/maps2"},
			want: want{
				statusCode: 400,
				location:   "",
			},
		},
		{
			name:    "absent url test #4",
			request: "/cc",
			storage: map[string]string{"aa": "https://gobyexample.com/maps", "bb": "https://gobyexample.com/maps2"},
			want: want{
				statusCode: 400,
				location:   "",
			},
		},
		{
			name:    "empty short url #5",
			request: "/",
			storage: map[string]string{"aa": "https://gobyexample.com/maps", "bb": "https://gobyexample.com/maps2"},
			want: want{
				statusCode: 400,
				location:   "",
			},
		},
		{
			name:    "empty storage #6",
			request: "/aa",
			storage: map[string]string{},
			want: want{
				statusCode: 400,
				location:   "",
			},
		},
		{
			name:    "url with a path #7",
			request: "/aa/bb",
			storage: map[string]string{"aa": "https://gobyexample.com/maps", "bb": "https://gobyexample.com/maps2"},
			want: want{
				statusCode: 400,
				location:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			request.SetPathValue("id", strings.TrimPrefix(tt.request, "/"))
			w := httptest.NewRecorder()
			repo := repository.NewUrlMemoryRepository(tt.storage)
			urlService := service.NewURLService(repo)
			h := http.Handler(NewURLGetHandler(urlService))
			h.ServeHTTP(w, request)

			result := w.Result()
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(result.Body)

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
		})
	}
}
