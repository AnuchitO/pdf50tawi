package pdf50tawi

import (
	"bytes"
	"image"
	"image/png"
	"testing"
)

func TestTinyEmptyPNG(t *testing.T) {
	// Test that tinyEmptyPNG returns valid PNG data
	pngData := tinyEmptyPNG()
	if pngData == nil {
		t.Fatal("expected non-nil PNG data")
	}

	if len(pngData) == 0 {
		t.Fatal("expected non-empty PNG data")
	}

	// Verify it's a valid PNG by decoding it
	img, err := png.Decode(bytes.NewReader(pngData))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	// Verify it's 1x1 pixels
	bounds := img.Bounds()
	if bounds.Dx() != 1 || bounds.Dy() != 1 {
		t.Fatalf("expected 1x1 image, got %dx%d", bounds.Dx(), bounds.Dy())
	}

	// Verify the pixel is transparent
	nrgbaImg, ok := img.(*image.NRGBA)
	if !ok {
		t.Fatalf("expected NRGBA image: %#v", nrgbaImg)
	}
	pixel := nrgbaImg.NRGBAAt(0, 0)
	if pixel.A != 0 {
		t.Fatalf("expected transparent pixel (A=0), got A=%d", pixel.A)
	}
}

func TestTinyEmptyPNGConsistency(t *testing.T) {
	// Test that multiple calls return identical data
	png1 := tinyEmptyPNG()
	png2 := tinyEmptyPNG()

	if len(png1) != len(png2) {
		t.Fatalf("expected consistent length: %d vs %d", len(png1), len(png2))
	}

	if !bytes.Equal(png1, png2) {
		t.Fatal("expected identical PNG data from multiple calls")
	}
}
