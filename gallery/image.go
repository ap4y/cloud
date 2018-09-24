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

// func init() {
// 	imagick.Initialize()
// }

// func ThumbnailForFile(file *os.File, maxHeight uint) (*os.File, error) {
// 	mw := imagick.NewMagickWand()
// 	defer mw.Destroy()

// 	if err := mw.ReadImageFile(file); err != nil {
// 		return nil, fmt.Errorf("failed to decode image: %s", err)
// 	}

// 	// log.Println(mw.GetImageProperties("exif:*"))
// 	width := (float32(maxHeight) / float32(mw.GetImageHeight())) * float32(mw.GetImageWidth())
// 	if err := mw.ThumbnailImage(uint(width), maxHeight); err != nil {
// 		return nil, fmt.Errorf("libmagick error: %s", err)
// 	}

// 	imgFile, err := ioutil.TempFile("", "thumb")
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create temp file: %s", err)
// 	}

// 	if err := mw.WriteImageFile(imgFile); err != nil {
// 		return nil, fmt.Errorf("failed to write file: %s", err)
// 	}

// 	return imgFile, nil
// }

func ThumbnailForFile(file *os.File, maxHeight uint) (io.Reader, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %s", err)
	}

	thumb := resize.Resize(0, maxHeight, img, resize.Lanczos3)
	out := bytes.NewBuffer([]byte{})
	if err := jpeg.Encode(out, thumb, nil); err != nil {
		return nil, fmt.Errorf("falied to encode jpeg: %s", err)
	}

	return out, nil
}

func EXIFForFile(file *os.File) (*exif.Exif, error) {
	x, err := exif.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode EXIF: %s", err)
	}

	return x, nil
}
