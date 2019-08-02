package gallery

import (
	"image"
	"os"
	"testing"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestThumbnail(t *testing.T) {
	file, err := os.Open("fixtures/album1/test.jpg")
	require.NoError(t, err)

	thumb, err := Thumbnail(file, 200)
	require.NoError(t, err)
	require.NotNil(t, thumb)

	cfg, _, err := image.DecodeConfig(thumb)
	require.NoError(t, err)
	assert.Equal(t, 200, cfg.Height)
}

func TestEXIF(t *testing.T) {
	file, err := os.Open("fixtures/album1/test.jpg")
	require.NoError(t, err)

	x, err := EXIF(file)
	require.NoError(t, err)
	require.NotNil(t, x)

	tag, err := x.Get(exif.PixelXDimension)
	require.NoError(t, err)
	assert.Equal(t, "2000", tag.String())
}
