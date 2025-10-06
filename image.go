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

// LoadOption represents an option for LoadImage function
type LoadOption func(*loadOptions)

type loadOptions struct {
	httpRequest *http.Request
}

// WithHTTPRequest adds HTTP request context for multipart file uploads
func WithHTTPRequest(r *http.Request) LoadOption {
	return func(opts *loadOptions) {
		opts.httpRequest = r
	}
}

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

// LoadImage loads an image based on its source type
// For Upload source type, pass WithHTTPRequest(r) option
// For File and URL source types, no options needed
func LoadImage(image Image, options ...LoadOption) (io.ReadCloser, error) {
	opts := &loadOptions{}
	for _, opt := range options {
		opt(opts)
	}

	switch image.SourceType {
	case Upload:
		if opts.httpRequest == nil {
			return nil, fmt.Errorf("LoadImage: Upload source type requires HTTP request context. Use LoadImage(image, WithHTTPRequest(r))")
		}
		return LoadImageFromMultiPartFile(opts.httpRequest, image.Value)
	case URL:
		return LoadImageFromURL(image.Value)
	case File:
		return LoadImageFromFile(image.Value)
	case "":
		return io.NopCloser(bytes.NewReader(tinyEmptyPNG())), nil
	default:
		return nil, fmt.Errorf("LoadImage: unsupported source type: %s", image.SourceType)
	}
}

// LoadImageFromFile loads image from file path (no HTTP context needed)
func LoadImageFromFile(file string) (io.ReadCloser, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, f)
	if err != nil {
		return nil, err
	}

	return io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

// LoadImageFromMultiPartFile loads image from multipart form upload (requires HTTP context)
// Use this function in web handlers when dealing with file uploads
func LoadImageFromMultiPartFile(r *http.Request, file string) (io.ReadCloser, error) {
	f, header, err := r.FormFile(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if header.Header.Get("Content-Type") != "image/png" {
		return nil, fmt.Errorf("invalid content type: %s", header.Header.Get("Content-Type"))
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, f)
	if err != nil {
		return nil, err
	}

	return io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

// LoadImageFromURL loads image from URL (no HTTP context needed)
func LoadImageFromURL(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "image/png" {
		return nil, fmt.Errorf("invalid content type: %s", resp.Header.Get("Content-Type"))
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return nil, err
	}

	return io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
