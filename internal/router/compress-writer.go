package router

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type compressWriter struct {
	w           http.ResponseWriter
	zw          *gzip.Writer
	wroteHeader bool
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w: w,
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	if c.zw == nil && !c.wroteHeader {
		c.startCompression(http.StatusOK)
	}
	if c.zw == nil {
		return c.w.Write(p)
	}
	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if c.wroteHeader {
		return
	}
	c.wroteHeader = true
	c.startCompression(statusCode)
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) startCompression(statusCode int) {
	if c.zw != nil || statusCode >= 300 || !c.contentZippable() {
		return
	}
	c.zw = gzip.NewWriter(c.w)
	c.w.Header().Set("Content-Encoding", "gzip")
}

func (c *compressWriter) contentZippable() bool {
	contentType := c.w.Header().Get("Content-Type")
	return strings.Contains(contentType, "text/html") || strings.Contains(contentType, "application/json")
}

func (c *compressWriter) Close() error {
	if c.zw == nil {
		return nil
	}
	return c.zw.Close()
}
