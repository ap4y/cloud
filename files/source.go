package files

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ItemType defines entity types supported by sources.
type ItemType string

const (
	// ItemTypeFile represents files.
	ItemTypeFile ItemType = "file"
	// ItemTypeDirectory represents directories.
	ItemTypeDirectory ItemType = "directory"
)

// Item represents a single source file.
type Item struct {
	Type    ItemType  `json:"type"`
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	ModTime time.Time `json:"updated_at"`
}

// Source provides album and images metadata.
type Source interface {
	List(path string) ([]*Item, error)
	Content(filePath string) (io.ReadCloser, error)
	Save(r io.Reader, path, filename string) error
	Remove(filePath string) error
}

type diskSource struct {
	basePath string
}

// NewDiskSource returns disk based source for a provided base dir
// path. Only images with requested extensions are returned.
func NewDiskSource(basePath string) (Source, error) {
	if !filepath.IsAbs(basePath) {
		return nil, errors.New("path is not absolute")
	}

	fi, err := os.Stat(basePath)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %s", err)
	}

	if !fi.IsDir() {
		return nil, errors.New("path is not directory")
	}

	return &diskSource{basePath}, nil
}

func (ds *diskSource) List(path string) ([]*Item, error) {
	cleanPath := strings.ReplaceAll(filepath.Clean(path), "..", "")
	diskPath := filepath.Join(ds.basePath, cleanPath)

	items := []*Item{}
	err := filepath.Walk(diskPath, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk %s: %s", path, err)
		}

		var item *Item
		if fi.IsDir() {
			item = &Item{
				Type:    ItemTypeDirectory,
				Name:    fi.Name(),
				Path:    fi.Name(),
				ModTime: fi.ModTime(),
			}
		} else {
			item = &Item{
				Type:    ItemTypeFile,
				Name:    strings.ReplaceAll(fi.Name(), filepath.Ext(fi.Name()), ""),
				Path:    strings.ReplaceAll(path, ds.basePath, ""),
				ModTime: fi.ModTime(),
			}
		}

		if path != diskPath {
			items = append(items, item)
		}

		if fi.IsDir() && path != diskPath {
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walk %s: %s", diskPath, err)
	}

	return items, nil
}

func (ds *diskSource) Content(filePath string) (io.ReadCloser, error) {
	cleanPath := strings.ReplaceAll(filepath.Clean(filePath), "..", "")
	diskPath := filepath.Join(ds.basePath, cleanPath)

	file, err := os.Open(diskPath)
	if err != nil {
		return nil, fmt.Errorf("open %s: %s", diskPath, err)
	}

	return file, nil
}

func (ds *diskSource) Save(r io.Reader, path, filename string) error {
	filePath := filepath.Join(path, strings.ReplaceAll(filename, "..", ""))
	cleanPath := strings.ReplaceAll(filepath.Clean(filePath), "..", "")
	diskPath := filepath.Join(ds.basePath, cleanPath)

	file, err := os.OpenFile(diskPath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return fmt.Errorf("file: %s", err)
	}

	if _, err := io.Copy(file, r); err != nil {
		return fmt.Errorf("file: %s", err)
	}

	return nil
}

func (ds *diskSource) Remove(filePath string) error {
	cleanPath := strings.ReplaceAll(filepath.Clean(filePath), "..", "")
	diskPath := filepath.Join(ds.basePath, cleanPath)

	return os.Remove(diskPath)
}
