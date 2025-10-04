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
	case SourceTypeUpload:
		return LoadImageFromMultiPartFile(r, image.Value)
	case SourceTypeURL:
		return LoadImageFromURL(image.Value)
	case SourceTypeFile:
		return os.Open(image.Value)
	default:
		return nil, fmt.Errorf("invalid source type: %s", image.SourceType)
	}
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

	return f, nil
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return io.NopCloser(bytes.NewReader(body)), nil
}
