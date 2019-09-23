package gallery

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"

	"github.com/nfnt/resize"
	"github.com/rwcarlsen/goexif/exif"
)

// Thumbnail returns thubmnail from a reader constrained by maxHeight.
func Thumbnail(r io.Reader, maxHeight uint) (io.Reader, error) {
	img, _, err := image.Decode(r)
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

// EXIF returns exif metadata from a reader.
func EXIF(r io.Reader) (*exif.Exif, error) {
	x, err := exif.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("decode: %s", err)
	}

	return x, nil
}
