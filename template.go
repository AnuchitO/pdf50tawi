package pdf50tawi

import (
	"bytes"
	"embed"
	"io"
)

//go:embed form
var form embed.FS

func certificateTemplate() (io.ReadSeeker, error) {
	f, err := form.ReadFile("form/tax50tawiTemplate.pdf")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(f), nil
}
