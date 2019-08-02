package gallery

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Cache caches gallery metadata.
type Cache interface {
	// Returns thumbnail and modtime for a given image path if it exists, otherwise returns nil.
	Thumbnail(imagePath string) (io.ReadSeeker, time.Time)
	// Stores thumbnail for a given image path and returns reader for a stored thumbnail.
	StoreThumbnail(imagePath string, r io.Reader) (io.ReadSeeker, error)
}

// DiskCache implements cache over filesystem.
type diskCache struct {
	dir string
}

// NewDiskCache returns a new Cache instance that uses filesystem.
// dir will be created if necessary.
func NewDiskCache(dir string) (Cache, error) {
	if dir == "" {
		return nil, errors.New("dir can't be empty")
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create cache dir: %s", err)
		}
	}

	return &diskCache{dir}, nil
}

func (dc *diskCache) cacheKey(imagePath string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(imagePath)))
}

func (dc *diskCache) Thumbnail(imagePath string) (io.ReadSeeker, time.Time) {
	filename := dc.cacheKey(imagePath)
	path := filepath.Join(dc.dir, filename)
	file, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		return nil, time.Time{}
	}

	fi, err := file.Stat()
	if err != nil {
		return file, time.Now()
	}

	return file, fi.ModTime()
}

func (dc *diskCache) StoreThumbnail(imagePath string, r io.Reader) (io.ReadSeeker, error) {
	filename := dc.cacheKey(imagePath)
	path := filepath.Join(dc.dir, filename)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return nil, fmt.Errorf("file: %s", err)
	}

	if _, err := io.Copy(file, r); err != nil {
		return nil, fmt.Errorf("copy: %s", err)
	}

	return file, nil
}
