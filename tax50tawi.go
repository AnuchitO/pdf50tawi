package pdf50tawi

import (
	"bytes"
	"embed"
	"io"
)

//go:embed file/tax50tawi.pdf
var tax50tawiPDF embed.FS

func Tax50tawiPDFTemplate() (io.ReadSeeker, error) {
	f, err := tax50tawiPDF.ReadFile("file/tax50tawi.pdf")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(f), nil
}
