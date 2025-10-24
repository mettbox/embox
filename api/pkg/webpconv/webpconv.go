package webpconv

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

type WebPResult struct {
	Data       []byte
	HasGeo     bool
	Latitude   float64
	Longitude  float64
	IsPortrait bool
}

// ConvertToWebP converts an image to WebP format, resizing it if necessary.
// It also extracts EXIF metadata for geolocation and orientation.
// If the image is already in WebP format, it returns the original data.
func ConvertToWebP(data []byte, maxSize int, quality float32) (*WebPResult, error) {
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	if format == "webp" {
		return &WebPResult{Data: data}, nil
	}

	var lat, lon float64
	var hasGeo bool

	exifData, err := exif.Decode(bytes.NewReader(data))
	if err == nil {
		lat, lon, err = exifData.LatLong()
		if err == nil {
			hasGeo = true
		}
		orientTag, err := exifData.Get(exif.Orientation)
		if err == nil {
			orient, err := orientTag.Int(0)
			if err == nil {
				img = applyOrientation(img, orient)
			}
		}
	}

	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	isPortrait := h > w

	if w > maxSize || h > maxSize {
		img = imaging.Fit(img, maxSize, maxSize, imaging.Lanczos)
	}

	var buf bytes.Buffer
	if err := webp.Encode(&buf, img, &webp.Options{Lossless: false, Quality: quality}); err != nil {
		return nil, fmt.Errorf("failed to encode webp: %w", err)
	}

	return &WebPResult{
		Data:       buf.Bytes(),
		Latitude:   lat,
		Longitude:  lon,
		HasGeo:     hasGeo,
		IsPortrait: isPortrait,
	}, nil
}

func applyOrientation(img image.Image, orient int) image.Image {
	switch orient {
	case 3:
		return imaging.Rotate180(img)
	case 6:
		return imaging.Rotate270(img)
	case 8:
		return imaging.Rotate90(img)
	default:
		return img
	}
}
