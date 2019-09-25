package gallery

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiskCache(t *testing.T) {
	dir, err := ioutil.TempDir("", "cache")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	cache, err := NewDiskCache(dir)
	require.NoError(t, err)

	img, _ := cache.Thumbnail("", "test.jpg")
	require.Nil(t, img)

	file, err := os.Open("fixtures/album1/test.jpg")
	require.NoError(t, err)

	ts := time.Now()
	_, err = cache.StoreThumbnail("", "test.jpg", file)
	require.NoError(t, err)

	img, modtime := cache.Thumbnail("", "test.jpg")
	require.NotNil(t, img)
	assert.Equal(t, ts.Unix(), modtime.Unix())
}
