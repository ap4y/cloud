package gallery

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"

	"github.com/nfnt/resize"
	"github.com/rwcarlsen/goexif/exif"
)

// Thumbnail returns thubmnail for an image file constrained by maxHeight.
func Thumbnail(file *os.File, maxHeight uint) (io.Reader, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("decode: %s", err)
	}

	thumb := resize.Resize(0, maxHeight, img, resize.Lanczos3)
	out := bytes.NewBuffer([]byte{})
	if err := jpeg.Encode(out, thumb, nil); err != nil {
		return nil, fmt.Errorf("encode: %s", err)
	}

	return out, nil
}

// EXIF returns exif metadata for an image file.
func EXIF(file *os.File) (*exif.Exif, error) {
	x, err := exif.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("decode: %s", err)
	}

	return x, nil
}
