package cache


import (
	"os"
	"go.uber.org/zap"
	"strings"
	"errors"
)

type Store struct {
	loc string
	ttl int
	logger *zap.Logger
}

const (
	CACHE_DIR = "sidekick-cache"
)


func NewStore(loc string, ttl int, logger *zap.Logger) *Store {
	os.MkdirAll(loc+"/"+CACHE_DIR, os.ModePerm)

	return &Store{
		loc: loc,
		ttl: ttl,
		logger: logger,
	}
}


func (d *Store) Get(key string) ([]byte, error) {
	key = strings.ReplaceAll(key, "/", "+")

	if _, err := os.Stat(d.loc + "/" + CACHE_DIR+ "/." +key); err == nil {
		d.logger.Debug("Pulled key from file")
		return os.ReadFile(d.loc+"/"+CACHE_DIR+"/."+key)
	}

	err := errors.New("Key not found in cache")

	return nil, err
}

func (d *Store) Set(key string, value []byte) error {
	key = strings.ReplaceAll(key, "/", "+")

    return os.WriteFile(d.loc+"/"+CACHE_DIR+"/."+key, value, 0644)
}

func (d *Store) Purge(key string) {
	key = strings.ReplaceAll(key, "/", "+")
	removeLoc := d.loc+"/"+CACHE_DIR+"/."
	d.logger.Debug("Removing key from cache", zap.String("key", key))
	
	if _, err := os.Stat(removeLoc+"br::"+key); err == nil {
		d.logger.Info("Removing brotli cache")
		err = os.Remove(removeLoc+"br::"+key)
		d.logger.Info("Brotli remove error status", zap.Error(err))
	}

	if _, err := os.Stat(removeLoc+"gzip::"+key); err == nil {
		d.logger.Info("Removing gzip cache")
		err = os.Remove(removeLoc+"gzip::"+key)
		d.logger.Info("Gzip remove error status", zap.Error(err))
	}
}

func (d *Store) Flush() error {
	return os.RemoveAll(d.loc + "/" + CACHE_DIR)
}