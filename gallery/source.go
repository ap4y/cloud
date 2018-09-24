package gallery

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Source interface {
	Albums() ([]*Album, error)
	Images(galleryName string) ([]*Image, error)
	Image(imagePath string) (*os.File, error)
}

type DiskSource struct {
	basePath      string
	imgExtensions map[string]bool
}

func NewDiskSource(basePath string, imgExtensions []string) (*DiskSource, error) {
	if !filepath.IsAbs(basePath) {
		return nil, fmt.Errorf("path is not absolute")
	}

	fi, err := os.Stat(basePath)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %s", err)
	}

	if !fi.IsDir() {
		return nil, fmt.Errorf("path is not directory")
	}

	exts := map[string]bool{}
	for _, ext := range imgExtensions {
		exts[strings.ToLower(ext)] = true
	}

	return &DiskSource{basePath, exts}, nil
}

func (ds *DiskSource) Albums() ([]*Album, error) {
	fis, err := ioutil.ReadDir(ds.basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to scan directory: %s", err)
	}

	albums := []*Album{}
	for _, fi := range fis {
		images, err := ds.images(fi.Name())
		if err != nil {
			continue
		}

		albums = append(albums, &Album{fi.Name(), fi.ModTime(), len(images)})
	}

	return albums, nil
}

func (ds *DiskSource) Images(galleryName string) ([]*Image, error) {
	cleanPath := strings.Replace(filepath.Clean(galleryName), "..", "", -1)
	images, err := ds.images(cleanPath)
	if err != nil {
		return nil, err
	}

	return images, nil
}

func (ds *DiskSource) Image(imagePath string) (*os.File, error) {
	cleanPath := strings.Replace(filepath.Clean(imagePath), "..", "", -1)
	diskPath := filepath.Join(ds.basePath, cleanPath)

	file, err := os.Open(diskPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %s", diskPath, err)
	}

	if !ds.imgExtensions[strings.ToLower(filepath.Ext(diskPath))] {
		return nil, fmt.Errorf("unknown file type: %s", diskPath)
	}

	return file, nil
}

func (ds *DiskSource) images(folderName string) ([]*Image, error) {
	diskPath := filepath.Join(ds.basePath, folderName)
	images := []*Image{}

	err := filepath.Walk(diskPath, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed accessing a path %q: %s\n", diskPath, err)
		}

		if fi.IsDir() && path != diskPath {
			return filepath.SkipDir
		}

		ext := filepath.Ext(path)
		if fi.IsDir() || !ds.imgExtensions[strings.ToLower(ext)] {
			return nil
		}

		relPath, _ := filepath.Rel(ds.basePath, path)
		images = append(images, &Image{strings.Replace(fi.Name(), ext, "", -1), relPath, fi.ModTime()})
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk gallery path %s: %s", diskPath, err)
	}

	return images, nil
}
