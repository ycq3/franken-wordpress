package cache

import (
	"net/http"
	"go.uber.org/zap"
)

func NewCustomWriter(rw http.ResponseWriter, r *http.Request, db *Store, logger *zap.Logger, path string) *CustomWriter {	
	nw := CustomWriter{rw, r, db, logger, path}
	return &nw
}

// CustomWriter handles the response and provide the way to cache the value
type CustomWriter struct {
	http.ResponseWriter
	*http.Request
	*Store
	*zap.Logger
	path string
}

func (r *CustomWriter) Header() http.Header {
	return r.ResponseWriter.Header()
}

// Write will write the response body
func (r *CustomWriter) Write(b []byte) (int, error) {
	r.Logger.Debug("Writing to cache", zap.String("path", r.path))
	// content encoding
	ct := r.Header().Get("Content-Encoding")
	r.Header().Set("X-WPEverywhere-Cache", "MISS")

	cacheKey := ct + "::" + r.path
	r.Logger.Debug("Cache Key", zap.String("Key", cacheKey))

	go r.Store.Set(cacheKey, b)

	return r.ResponseWriter.Write(b)
}
