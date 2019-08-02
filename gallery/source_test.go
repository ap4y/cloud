package gallery

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSource(t *testing.T) {
	pwd, err := os.Getwd()
	require.NoError(t, err)

	src, err := NewDiskSource(filepath.Join(pwd, "fixtures"), []string{".jpg"})
	require.NoError(t, err)

	t.Run("Albums", func(t *testing.T) {
		albums, err := src.Albums()
		require.NoError(t, err)
		require.Len(t, albums, 2)

		album := albums[0]
		assert.Equal(t, "album1", album.Name)
		assert.Equal(t, 1, album.ItemsCount)
	})

	t.Run("Images", func(t *testing.T) {
		imgs, err := src.Images("album1")
		require.NoError(t, err)
		require.Len(t, imgs, 1)

		img := imgs[0]
		assert.Equal(t, "test", img.Name)
	})

	t.Run("Image", func(t *testing.T) {
		img, err := src.Image("/album1/test.jpg")
		require.NoError(t, err)
		assert.NotNil(t, img)

		img, err = src.Image("../album1/test.jpg")
		require.NoError(t, err)
		assert.NotNil(t, img)
	})
}

func TestSourceIgnoreExts(t *testing.T) {
	pwd, err := os.Getwd()
	require.NoError(t, err)

	src, err := NewDiskSource(filepath.Join(pwd, "fixtures"), []string{})
	require.NoError(t, err)

	t.Run("Albums", func(t *testing.T) {
		albums, err := src.Albums()
		require.NoError(t, err)
		require.Len(t, albums, 2)

		album := albums[0]
		assert.Equal(t, "album1", album.Name)
		assert.Equal(t, 0, album.ItemsCount)
	})

	t.Run("Images", func(t *testing.T) {
		imgs, err := src.Images("album1")
		require.NoError(t, err)
		require.Len(t, imgs, 0)
	})

	t.Run("Image", func(t *testing.T) {
		_, err := src.Image("/album1/test.jpg")
		require.Error(t, err)
	})
}
