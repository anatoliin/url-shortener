package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anatoliin/url-shortener/internal/repository"
	"github.com/anatoliin/url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
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
				statusCode: 404,
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
				statusCode: 404,
				location:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repository.NewURLMemoryRepository(tt.storage)
			urlService := service.NewURLService(repo)

			r := chi.NewRouter()
			r.Method(http.MethodGet, "/{id}", NewURLGetHandler(urlService))
			ts := httptest.NewServer(r)
			defer ts.Close()

			resp, _ := testRequest(t, ts, "GET", tt.request, nil)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.location, resp.Header.Get("Location"))
		})
	}
}
