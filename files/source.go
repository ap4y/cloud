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
	Type     ItemType  `json:"type"`
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	ModTime  time.Time `json:"updated_at"`
	Children []*Item   `json:"children"`
}

// Source provides album and images metadata.
type Source interface {
	Tree() (*Item, error)
	File(filePath string) (*os.File, error)
	Save(r io.Reader, filePath string) (*Item, error)
	Remove(filePath string) (*Item, error)
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

func (ds *diskSource) Tree() (*Item, error) {
	items := map[string]*Item{}
	err := filepath.Walk(ds.basePath, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk %s: %s", path, err)
		}

		relPath := strings.ReplaceAll(path, ds.basePath, "")
		dirPath := strings.ReplaceAll(relPath, "/"+fi.Name(), "")

		var item *Item
		if fi.IsDir() {
			item = &Item{
				Type:     ItemTypeDirectory,
				Name:     fi.Name(),
				Path:     relPath,
				ModTime:  fi.ModTime(),
				Children: make([]*Item, 0),
			}
			items[relPath] = item
		} else {
			item = &Item{
				Type:    ItemTypeFile,
				Name:    fi.Name(),
				Path:    relPath,
				ModTime: fi.ModTime(),
			}
		}

		if path != ds.basePath {
			items[dirPath].Children = append(items[dirPath].Children, item)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walk: %s", err)
	}

	return &Item{
		Type:     ItemTypeDirectory,
		Name:     "/",
		Path:     "/",
		ModTime:  time.Now(),
		Children: (items[""].Children),
	}, nil
}

func (ds *diskSource) File(filePath string) (*os.File, error) {
	cleanPath := strings.ReplaceAll(filepath.Clean(filePath), "..", "")
	diskPath := filepath.Join(ds.basePath, cleanPath)

	file, err := os.Open(diskPath)
	if err != nil {
		return nil, fmt.Errorf("open %s: %s", diskPath, err)
	}

	return file, nil
}

func (ds *diskSource) Save(r io.Reader, filePath string) (*Item, error) {
	cleanPath := filepath.Clean(strings.ReplaceAll(filePath, "..", ""))
	diskPath := filepath.Join(ds.basePath, cleanPath)

	file, err := os.OpenFile(diskPath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, fmt.Errorf("file: %s", err)
	}

	if _, err := io.Copy(file, r); err != nil {
		return nil, fmt.Errorf("file: %s", err)
	}

	_, filename := filepath.Split(filePath)
	return &Item{
		Type:    ItemTypeFile,
		Name:    filename,
		Path:    cleanPath,
		ModTime: time.Now(),
	}, nil
}

func (ds *diskSource) Remove(filePath string) (*Item, error) {
	cleanPath := strings.ReplaceAll(filepath.Clean(filePath), "..", "")
	diskPath := filepath.Join(ds.basePath, cleanPath)

	if err := os.Remove(diskPath); err != nil {
		return nil, err
	}

	_, filename := filepath.Split(filePath)
	return &Item{
		Type:    ItemTypeFile,
		Name:    filename,
		Path:    cleanPath,
		ModTime: time.Now(),
	}, nil
}
