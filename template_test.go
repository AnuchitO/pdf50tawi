package pdf50tawi

import "testing"

func TestCertificateTemplate(t *testing.T) {
	r, err := certificateTemplate()
	if err != nil {
		t.Fatalf("certificateTemplate error: %v", err)
	}
	buf := make([]byte, 16)
	n, err := r.Read(buf)
	if err != nil && err.Error() != "EOF" && n == 0 {
		t.Fatalf("unable to read template: %v", err)
	}
}
