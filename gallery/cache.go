package gallery

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Cache interface {
	Thumbnail(imagePath string) (io.ReadSeeker, time.Time)
	StoreThumbnail(imagePath string, r io.Reader) (io.ReadSeeker, error)
}

type DiskCache struct {
	dir string
}

func NewDiskCache(dir string) (*DiskCache, error) {
	if dir == "" {
		dir = os.TempDir()
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create cache dir: %s", err)
		}
	}

	return &DiskCache{dir}, nil
}

func (dc *DiskCache) Thumbnail(imagePath string) (io.ReadSeeker, time.Time) {
	filename := fmt.Sprintf("%x", md5.Sum([]byte(imagePath)))
	name := filepath.Join(dc.dir, filename)
	file, err := os.OpenFile(name, os.O_RDONLY, 0600)
	if err != nil {
		return nil, time.Time{}
	}

	fi, err := file.Stat()
	if err != nil {
		return file, time.Now()
	}

	return file, fi.ModTime()
}

func (dc *DiskCache) StoreThumbnail(imagePath string, r io.Reader) (io.ReadSeeker, error) {
	filename := fmt.Sprintf("%x", md5.Sum([]byte(imagePath)))
	name := filepath.Join(dc.dir, filename)
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %s", err)
	}

	if _, err := io.Copy(file, r); err != nil {
		return nil, fmt.Errorf("failed to write to a file: %s", err)
	}

	return file, nil
}

func (dc *DiskCache) openThumbFile(filename string) (*os.File, error) {
	name := filepath.Join(dc.dir, filename)
	return os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
}
