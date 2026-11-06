package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	var output bytes.Buffer
	logger := zerolog.New(&output)

	handler := Logger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte("created"))
		require.NoError(t, err)
		time.Sleep(20 * time.Millisecond)
	}))
	request := httptest.NewRequest(http.MethodPost, "/short", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	require.Equal(t, http.StatusCreated, response.Code)
	require.Equal(t, "created", response.Body.String())

	var event map[string]any
	require.NoError(t, json.Unmarshal(output.Bytes(), &event))
	require.Equal(t, "/short", event["uri"])
	require.Equal(t, http.MethodPost, event["method"])
	require.Equal(t, float64(http.StatusCreated), event["status_code"])
	require.Equal(t, float64(len("created")), event["bytes_out"])
	duration, ok := event["duration"].(float64)
	require.True(t, ok, "duration should be logged as a number")
	require.GreaterOrEqual(t, duration, float64(19))
}
