package fileutil

import (
	"io/ioutil"
	"sync"
)

type readFileCacheItem struct {
	s string
	e error
}
type readFileCacheType struct {
	lock   sync.RWMutex
	byname map[string]readFileCacheItem
}

var readFileCache readFileCacheType

func init() {
	readFileCache.byname = make(map[string]readFileCacheItem)
}

// ReadFileNoCache Read a file from disk, return as a string.
func ReadFileNoCache(fn string) (string, error) {
	b, e := ioutil.ReadFile(fn)
	return string(b), e
}

// ReadFile will check the cache first, then fallback to ReadFileFromDisk
func ReadFile(fn string) (string, error) {
	readFileCache.lock.Lock()
	defer readFileCache.lock.Unlock()
	if item, ok := readFileCache.byname[fn]; ok {
		return item.s, item.e
	}

	// Crap. Go read it for real.
	s, e := ReadFileNoCache(fn)
	readFileCache.byname[fn] = readFileCacheItem{s: s, e: e}
	return s, e
}
