package pdf50tawi

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

func LoadImage(image Image, r *http.Request) (io.ReadCloser, error) {
	switch image.SourceType {
	case Upload:
		return LoadImageFromMultiPartFile(r, image.Value)
	case URL:
		return LoadImageFromURL(image.Value)
	case File:
		return LoadImageFromFile(image.Value)
	default:
		return nil, fmt.Errorf("invalid source type: %s", image.SourceType)
	}
}

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
