package cache


import (
	"os"
	"go.uber.org/zap"
	"strings"
	"errors"
	"time"
)

type Store struct {
	loc string
	ttl int
	logger *zap.Logger
	memCache map[string]*MemCacheItem
}

type MemCacheItem struct {
	value []byte
	timestamp int64
}

const (
	CACHE_DIR = "sidekick-cache"
)


func NewStore(loc string, ttl int, logger *zap.Logger) *Store {
	os.MkdirAll(loc+"/"+CACHE_DIR, os.ModePerm)
	memCache := make(map[string]*MemCacheItem)

	// Load cache from disk
	files, err := os.ReadDir(loc+"/"+CACHE_DIR)
	if err == nil {
		for _, file := range files {
			if !file.IsDir() {
				value, err := os.ReadFile(loc+"/"+CACHE_DIR+"/"+file.Name())
				stat, err := os.Stat(loc+"/"+CACHE_DIR+"/"+file.Name())

				if err == nil {
					memCache[file.Name()] = &MemCacheItem{
						value: value,
						timestamp: stat.ModTime().Unix(),
					}
				}
			}
		}
	}

	return &Store{
		loc: loc,
		ttl: ttl,
		logger: logger,
		memCache: memCache,
	}
}


func (d *Store) Get(key string) ([]byte, error) {
	key = strings.ReplaceAll(key, "/", "+")

	if d.memCache[key] != nil {
		d.logger.Debug("Pulled key from memory", zap.String("key", key))
		return d.memCache[key].value, nil
	}

	if _, err := os.Stat(d.loc + "/" + CACHE_DIR+ "/." +key); err == nil {
		d.logger.Debug("Pulled key from file", zap.String("key", key))
		return os.ReadFile(d.loc+"/"+CACHE_DIR+"/."+key)
	}

	err := errors.New("Key not found in cache")

	return nil, err
}

func (d *Store) Set(key string, value []byte) error {
	key = strings.ReplaceAll(key, "/", "+")
	d.memCache[key] = &MemCacheItem{
		value: value,
		timestamp: time.Now().Unix(),
	}

    return os.WriteFile(d.loc+"/"+CACHE_DIR+"/."+key, value, 0644)
}

func (d *Store) Purge(key string) {
	key = strings.ReplaceAll(key, "/", "+")
	removeLoc := d.loc+"/"+CACHE_DIR+"/."
	d.logger.Debug("Removing key from cache", zap.String("key", key))

	delete(d.memCache, "br::"+key)
	delete(d.memCache, "gzip::"+key)
	
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
	d.memCache = make(map[string]*MemCacheItem)
	return os.RemoveAll(d.loc + "/" + CACHE_DIR)
}

func (d *Store) List() map[string][]string {
	list := make(map[string][]string)
	list["mem"] = make([]string, len(d.memCache))
	memIdx := 0

	for key, _ := range d.memCache {
		list["mem"][memIdx] = key
		memIdx++
	}

	files, err := os.ReadDir(d.loc+"/"+CACHE_DIR)
	list["disk"] = make([]string, 0)

	if err == nil {
		for _, file := range files {
			if !file.IsDir() {
				list["disk"] = append(list["disk"], file.Name())
			}
		}
	}

	return list
}