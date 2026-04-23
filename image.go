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

// emptyPNG is a 1×1 transparent PNG computed once at package init.
// Returned by tinyEmptyPNG() to avoid re-encoding on every nil-image call.
var emptyPNG = func() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.Transparent)
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}()

func tinyEmptyPNG() []byte { return emptyPNG }

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

// LoadImageFromRequest fetches a PNG image by executing the given HTTP request.
// Use this when the image URL requires custom headers such as authentication:
//
//	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
//	req.Header.Set("Authorization", "Bearer "+token)
//	img, err := pdf50tawi.LoadImageFromRequest(req)
func LoadImageFromRequest(req *http.Request) (io.Reader, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d fetching %s", resp.StatusCode, req.URL)
	}

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, resp.Body); err != nil {
		return nil, err
	}
	return bytes.NewReader(buf.Bytes()), nil
}
