package gallery

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/ap4y/cloud/internal/pathutil"
)

// Source provides album and images metadata.
type Source interface {
	// Albums returns list of all albums in a Source.
	Albums() ([]Album, error)
	// Images returns images metadata for a given album.
	Images(album string) ([]Image, error)
	// Image returns image file for a given image path.
	Image(ablum, image string) (*os.File, error)
}

type diskSource struct {
	basePath      string
	imgExtensions map[string]bool
}

// NewDiskSource returns disk based source for a provided base dir
// path. Only images with requested extensions are returned.
func NewDiskSource(basePath string, imgExtensions []string) (Source, error) {
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

	exts := map[string]bool{}
	for _, ext := range imgExtensions {
		exts[strings.ToLower(ext)] = true
	}

	return &diskSource{basePath, exts}, nil
}

func (ds *diskSource) Albums() ([]Album, error) {
	fis, err := ioutil.ReadDir(ds.basePath)
	if err != nil {
		return nil, fmt.Errorf("scan: %s", err)
	}

	albums := make([]Album, 0, len(fis))
	for _, fi := range fis {
		images, err := ds.images(fi.Name())
		if err != nil {
			continue
		}

		albums = append(albums, Album{fi.Name(), fi.ModTime(), len(images)})
	}

	return albums, nil
}

func (ds *diskSource) Images(album string) ([]Image, error) {
	images, err := ds.images(album)
	if err != nil {
		return nil, err
	}

	return images, nil
}

func (ds *diskSource) Image(album, image string) (*os.File, error) {
	diskPath := pathutil.Join(ds.basePath, album, image)

	file, err := os.Open(diskPath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %s", diskPath, err)
	}

	if !ds.imgExtensions[strings.ToLower(filepath.Ext(diskPath))] {
		return nil, fmt.Errorf("unknown file: %s", diskPath)
	}

	return file, nil
}

func (ds *diskSource) images(folderName string) ([]Image, error) {
	diskPath := pathutil.Join(ds.basePath, folderName)
	images := make([]Image, 0)

	err := filepath.Walk(diskPath, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk %s: %s", diskPath, err)
		}

		if fi.IsDir() && path != diskPath {
			return filepath.SkipDir
		}

		ext := filepath.Ext(path)
		if fi.IsDir() || !ds.imgExtensions[strings.ToLower(ext)] {
			return nil
		}

		image := Image{
			Name:    strings.Replace(fi.Name(), ext, "", -1),
			Path:    fi.Name(),
			ModTime: fi.ModTime(),
		}
		images = append(images, image)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walk %s: %s", diskPath, err)
	}

	return images, nil
}
