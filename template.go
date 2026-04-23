package pdf50tawi

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"os"
	"sync"
)

//go:embed form
var form embed.FS

// tplBytesOnce caches the raw template PDF bytes read from the embedded FS.
var (
	tplBytesOnce sync.Once
	tplBytes     []byte
	tplBytesErr  error
)

// certificateTemplate returns a ReadSeeker over the embedded PDF template.
// The bytes are read from the embedded FS exactly once; subsequent calls
// return a new reader over the same cached slice.
func certificateTemplate() (io.ReadSeeker, error) {
	tplBytesOnce.Do(func() {
		b, err := form.ReadFile("form/tax50tawiTemplate.pdf")
		if err != nil {
			tplBytesErr = err
			return
		}
		tplBytes = b
	})
	if tplBytesErr != nil {
		return nil, tplBytesErr
	}
	return bytes.NewReader(tplBytes), nil
}

// tplFileOnce caches a temp file that gopdf.ImportPage can read from.
// The file is written once and reused for the process lifetime — it is never
// removed so concurrent calls to fillCertificate are safe.
var (
	tplFileOnce sync.Once
	tplFilePath string
	tplFileErr  error
)

func cachedTemplatePath() (string, error) {
	tplFileOnce.Do(func() {
		r, err := certificateTemplate()
		if err != nil {
			tplFileErr = err
			return
		}
		data, err := io.ReadAll(r)
		if err != nil {
			tplFileErr = fmt.Errorf("read template: %w", err)
			return
		}
		f, err := os.CreateTemp("", "pdf50tawi-tpl-*.pdf")
		if err != nil {
			tplFileErr = fmt.Errorf("create temp: %w", err)
			return
		}
		if _, err := f.Write(data); err != nil {
			f.Close()
			os.Remove(f.Name())
			tplFileErr = fmt.Errorf("write temp template: %w", err)
			return
		}
		f.Close()
		tplFilePath = f.Name()
		// Intentionally NOT removed — reused for the process lifetime.
	})
	return tplFilePath, tplFileErr
}
