package pdf50tawi

import (
	"bytes"
	"testing"
)

func TestIssueWHTCertificatePDF(t *testing.T) {
	png := tinyEmptyPNG()
	var out bytes.Buffer
	err := IssueWHTCertificatePDF(&out, sampleTaxInfo(), bytes.NewReader(png), bytes.NewReader(png))
	if err != nil {
		t.Fatalf("IssueWHTCertificatePDF error: %v", err)
	}
	if out.Len() == 0 {
		t.Fatal("expected non-empty PDF output")
	}
	// Valid PDF starts with %PDF
	if !bytes.HasPrefix(out.Bytes(), []byte("%PDF")) {
		t.Fatalf("output does not look like a PDF (first bytes: %q)", out.Bytes()[:min(8, out.Len())])
	}
	// Should be well under 1MB
	if out.Len() > 1024*1024 {
		t.Fatalf("output too large: %d bytes (expected < 1MB)", out.Len())
	}
}

func TestFillCertificateWithEmptyFields(t *testing.T) {
	var out bytes.Buffer
	if err := fillCertificate(nil, nil, &out); err != nil {
		t.Fatalf("fillCertificate error: %v", err)
	}
	if !bytes.HasPrefix(out.Bytes(), []byte("%PDF")) {
		t.Fatal("output does not look like a PDF")
	}
}

func TestAnchorToXY(t *testing.T) {
	cases := []struct {
		anchor    Anchor
		dx, dy    float64
		wantX, wantY float64
	}{
		{TopLeft, 58, -98, 58, 98},
		{TopCenter, 0, -50, pageWidth / 2, 50},
		{BottomRight, -109.5, 530, pageWidth - 109.5, pageHeight - 530},
		{BottomCenter, 69, 530, pageWidth/2 + 69, pageHeight - 530},
		{BottomLeft, 10, 20, 10, pageHeight - 20},
		{Center, 0, 0, pageWidth / 2, pageHeight / 2},
	}
	for _, c := range cases {
		x, y := anchorToXY(c.anchor, c.dx, c.dy)
		if abs(x-c.wantX) > 0.01 || abs(y-c.wantY) > 0.01 {
			t.Errorf("anchorToXY(%v, %v, %v) = (%v, %v), want (%v, %v)",
				c.anchor, c.dx, c.dy, x, y, c.wantX, c.wantY)
		}
	}
}

func abs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
