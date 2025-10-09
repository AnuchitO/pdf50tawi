package pdf50tawi

import (
	"bytes"
	"io"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func TestBuildWriteAndWHTCertificatePDF(t *testing.T) {
	png := tinyEmptyPNG()
	// BuildStampedContext
	texts := []TextStamp{{Text: "t", Dx: 10, Dy: -10, FontSize: 12, Position: types.TopLeft}}
	images := []ImageStamp{{Reader: bytes.NewReader(png), Pos: types.BottomLeft, Dx: 5, Dy: 5, Scale: 0.1, Opacity: 1, OnTop: true}}
	ctx, err := BuildStampedContext(texts, images)
	if err != nil || ctx == nil {
		t.Fatalf("BuildStampedContext error: %v", err)
	}
	// WriteStampedPDF
	var out bytes.Buffer
	if err := WriteStampedPDF(ctx, &out); err != nil {
		t.Fatalf("WriteStampedPDF error: %v", err)
	}
	if out.Len() == 0 {
		t.Fatalf("expected output PDF bytes")
	}
	// WHTCertificatePDF using embedded template
	var out2 bytes.Buffer
	if err := IssueWHTCertificatePDF(&out2, sampleTaxInfo(), bytes.NewReader(png), bytes.NewReader(png)); err != nil {
		t.Fatalf("WHTCertificatePDF error: %v", err)
	}
	if out2.Len() == 0 {
		t.Fatalf("expected output bytes for WHTCertificatePDF")
	}
}

func TestTextWatermark(t *testing.T) {
	t.Run("default font", func(t *testing.T) {
		wm, err := TextWatermark(TextStamp{Text: "X", FontSize: 12})
		if err != nil {
			t.Fatal(err)
		}
		if wm.FontName != "THSarabunNew" {
			t.Fatalf("expected default font THSarabunNew, got %s", wm.FontName)
		}
	})
	t.Run("custom font", func(t *testing.T) {
		cfg := TextStamp{Text: "Hello", Dx: 10, Dy: 20, FontSize: 16, FontName: "CustomFont", Position: types.TopLeft}
		wm, err := TextWatermark(cfg)
		if err != nil {
			t.Fatalf("TextWatermark error: %v", err)
		}
		if wm.Dx != cfg.Dx || wm.Dy != cfg.Dy || wm.FontSize != cfg.FontSize || wm.FontName != cfg.FontName || wm.Pos != cfg.Position || !wm.ScaleAbs || wm.OnTop != true {
			t.Fatalf("unexpected watermark fields: %+v", wm)
		}
		if wm.FontName != "CustomFont" {
			t.Fatalf("expected font name CustomFont, got %s", wm.FontName)
		}
	})

	t.Run("empty text should return ' ' one space otherwise it will crash", func(t *testing.T) {
		wm, err := TextWatermark(TextStamp{Text: "", FontSize: 12})
		if err != nil {
			t.Fatal(err)
		}
		if wm.TextString != " " {
			t.Fatalf("expected text ' ', got %s", wm.TextString)
		}
	})

}

func TestImageWatermark(t *testing.T) {
	img := tinyEmptyPNG()
	wm, err := ImageWatermark(ImageStamp{Reader: bytes.NewReader(img), Pos: types.BottomRight, Dx: 3, Dy: 4, Scale: 0.5, Opacity: 0.8, OnTop: true})
	if err != nil {
		t.Fatalf("ImageWatermark error: %v", err)
	}
	if wm.Dx != 3 || wm.Dy != 4 || wm.Scale != 0.5 || !wm.ScaleAbs || wm.Opacity != 0.8 || wm.Pos != types.BottomRight || wm.OnTop != true {
		t.Fatalf("unexpected image watermark fields: %+v", wm)
	}
}

func mustPDFTemplate(t *testing.T) io.ReadSeeker {
	t.Helper()
	r, err := Tax50tawiPDFTemplate()
	if err != nil {
		t.Fatalf("Tax50tawiPDFTemplate error: %v", err)
	}
	return r
}

func TestReadContext(t *testing.T) {
	ctx, err := ReadContext(mustPDFTemplate(t))
	if err != nil || ctx == nil {
		t.Fatalf("ReadContext failed: %v, ctx=%v", err, ctx)
	}
	if ctx.Conf == nil {
		t.Fatalf("expected configuration present")
	}
}
