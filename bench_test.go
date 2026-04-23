package pdf50tawi

import (
	"bytes"
	"io"
	"testing"
)

// ── shared fixtures ──────────────────────────────────────────────────────────

func benchTaxInfo() TaxInfo { return sampleTaxInfo() }

// ── end-to-end ───────────────────────────────────────────────────────────────

// BenchmarkIssueWHTCertificatePDF measures the full PDF generation path
// including template loading, font embedding, text placement, and PDF write.
func BenchmarkIssueWHTCertificatePDF(b *testing.B) {
	tax := benchTaxInfo()
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if err := IssueWHTCertificatePDF(io.Discard, tax, nil, nil); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkIssueWHTCertificatePDFWithImages adds two real PNG images to the
// hot path to measure the extra cost of temp-file roundtrip per image.
func BenchmarkIssueWHTCertificatePDFWithImages(b *testing.B) {
	tax := benchTaxInfo()
	png := tinyEmptyPNG()
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		sign := bytes.NewReader(png)
		seal := bytes.NewReader(png)
		if err := IssueWHTCertificatePDF(io.Discard, tax, sign, seal); err != nil {
			b.Fatal(err)
		}
	}
}

// ── sub-components ───────────────────────────────────────────────────────────

// BenchmarkTextFieldsFromTaxInfo isolates the field-mapping step.
func BenchmarkTextFieldsFromTaxInfo(b *testing.B) {
	tax := benchTaxInfo()
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_ = TextFieldsFromTaxInfo(tax)
	}
}

// BenchmarkCertificateImageFields isolates image-field construction.
func BenchmarkCertificateImageFields(b *testing.B) {
	png := tinyEmptyPNG()
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_ = CertificateImageFields(bytes.NewReader(png), bytes.NewReader(png))
	}
}

// BenchmarkFillCertificateTextOnly measures template load + font + text
// placement without any images.
func BenchmarkFillCertificateTextOnly(b *testing.B) {
	tax := benchTaxInfo()
	texts := TextFieldsFromTaxInfo(tax)
	images := CertificateImageFields(nil, nil)
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if err := fillCertificate(texts, images, io.Discard); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkFillCertificateWithImages adds image placement to the above.
func BenchmarkFillCertificateWithImages(b *testing.B) {
	tax := benchTaxInfo()
	texts := TextFieldsFromTaxInfo(tax)
	png := tinyEmptyPNG()
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		images := CertificateImageFields(bytes.NewReader(png), bytes.NewReader(png))
		if err := fillCertificate(texts, images, io.Discard); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCertificateTemplate measures the cost of loading the embedded
// template on each call (the pre-cache baseline).
func BenchmarkCertificateTemplate(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		r, err := certificateTemplate()
		if err != nil {
			b.Fatal(err)
		}
		if _, err := io.Copy(io.Discard, r); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkTinyEmptyPNG measures the cost of generating the fallback PNG.
func BenchmarkTinyEmptyPNG(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_ = tinyEmptyPNG()
	}
}
