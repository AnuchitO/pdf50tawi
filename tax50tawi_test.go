package pdf50tawi

import "testing"

func TestTax50tawiPDFTemplate(t *testing.T) {
	r, err := Tax50tawiPDFTemplate()
	if err != nil {
		t.Fatalf("Tax50tawiPDFTemplate error: %v", err)
	}
	buf := make([]byte, 16)
	n, err := r.Read(buf)
	if err != nil && err.Error() != "EOF" && n == 0 {
		t.Fatalf("unable to read template: %v", err)
	}
}
