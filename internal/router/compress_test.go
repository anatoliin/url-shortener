package router

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGzipMiddleware(t *testing.T) {
	const requestBody = "compressed request"
	var compressedRequest bytes.Buffer
	writer := gzip.NewWriter(&compressedRequest)
	_, err := writer.Write([]byte(requestBody))
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	tests := []struct {
		name                string
		acceptEncoding      string
		contentEncoding     string
		requestBody         []byte
		responseStatus      int
		responseContentType string
		responseBody        string
		wantStatus          int
		wantEncoding        string
		wantResponseBody    string
		wantRequestBody     string
		wantHandlerCall     bool
		wantResponseGzip    bool
	}{
		{
			name:                "compresses accepted response",
			acceptEncoding:      "gzip",
			responseStatus:      http.StatusCreated,
			responseContentType: "text/html",
			responseBody:        "compressed response",
			wantStatus:          http.StatusCreated,
			wantEncoding:        "gzip",
			wantResponseBody:    "compressed response",
			wantHandlerCall:     true,
			wantResponseGzip:    true,
		},
		{
			name:                "does not append gzip data to plain response",
			acceptEncoding:      "gzip",
			responseStatus:      http.StatusCreated,
			responseContentType: "text/plain",
			responseBody:        "http://localhost:8080/kJJrauYx",
			wantStatus:          http.StatusCreated,
			wantResponseBody:    "http://localhost:8080/kJJrauYx",
			wantHandlerCall:     true,
		},
		{
			name:             "does not compress unaccepted response",
			responseStatus:   http.StatusOK,
			responseBody:     "plain response",
			wantStatus:       http.StatusOK,
			wantResponseBody: "plain response",
			wantHandlerCall:  true,
		},
		{
			name:            "decompresses request",
			contentEncoding: "gzip",
			requestBody:     compressedRequest.Bytes(),
			responseStatus:  http.StatusOK,
			wantStatus:      http.StatusOK,
			wantRequestBody: requestBody,
			wantHandlerCall: true,
		},
		{
			name:            "rejects invalid gzip request",
			contentEncoding: "gzip",
			requestBody:     []byte("not gzip"),
			wantStatus:      http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerCalled := false
			var receivedRequestBody string
			handler := GzipMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				receivedRequestBody = string(body)
				if tt.responseContentType != "" {
					w.Header().Set("Content-Type", tt.responseContentType)
				}
				w.WriteHeader(tt.responseStatus)
				_, err = w.Write([]byte(tt.responseBody))
				require.NoError(t, err)
			}))

			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(tt.requestBody))
			request.Header.Set("Accept-Encoding", tt.acceptEncoding)
			request.Header.Set("Content-Encoding", tt.contentEncoding)
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			require.Equal(t, tt.wantStatus, response.Code)
			require.Equal(t, tt.wantEncoding, response.Header().Get("Content-Encoding"))
			require.Equal(t, tt.wantHandlerCall, handlerCalled)
			require.Equal(t, tt.wantRequestBody, receivedRequestBody)

			if tt.wantResponseGzip {
				reader, err := gzip.NewReader(response.Body)
				require.NoError(t, err)
				defer reader.Close()
				body, err := io.ReadAll(reader)
				require.NoError(t, err)
				require.Equal(t, tt.wantResponseBody, string(body))
				return
			}

			if tt.wantHandlerCall {
				require.Equal(t, tt.wantResponseBody, response.Body.String())
			}
		})
	}
}
