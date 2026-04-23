package pdf50tawi

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"
	"os"
)

func tinyEmptyPNG() []byte {
	size := 1
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := range size {
		for x := range size {
			img.Set(x, y, color.Transparent)
		}
	}
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

// LoadImageFromFile loads a PNG or JPEG image from a local file path.
func LoadImageFromFile(file string) (io.Reader, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, f); err != nil {
		return nil, err
	}
	return bytes.NewReader(buf.Bytes()), nil
}

// LoadImageFromMultiPartFile reads a PNG image uploaded via multipart form.
// field is the form field name (e.g. "signature").
func LoadImageFromMultiPartFile(r *http.Request, field string) (io.Reader, error) {
	f, header, err := r.FormFile(field)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if header.Header.Get("Content-Type") != "image/png" {
		return nil, fmt.Errorf("invalid content type: %s", header.Header.Get("Content-Type"))
	}

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, f); err != nil {
		return nil, err
	}
	return bytes.NewReader(buf.Bytes()), nil
}

// LoadImageFromURL fetches a PNG image from the given URL.
func LoadImageFromURL(url string) (io.Reader, error) {
	resp, err := http.Get(url) //nolint:noctx
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d fetching %s", resp.StatusCode, url)
	}

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, resp.Body); err != nil {
		return nil, err
	}
	return bytes.NewReader(buf.Bytes()), nil
}
