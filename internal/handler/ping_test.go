package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type stubPinger struct {
	err error
}

func (p stubPinger) Ping(context.Context) error {
	return p.err
}

func TestPingHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name       string
		pingErr    error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "successful ping",
			wantStatus: http.StatusOK,
		},
		{
			name:       "failed ping",
			pingErr:    errors.New("database unavailable"),
			wantStatus: http.StatusInternalServerError,
			wantBody:   "database unavailable\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewPingHandler(stubPinger{err: tt.pingErr})
			request := httptest.NewRequest(http.MethodGet, "/ping", nil)
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			assert.Equal(t, tt.wantStatus, response.Code)
			assert.Equal(t, tt.wantBody, response.Body.String())
		})
	}
}
