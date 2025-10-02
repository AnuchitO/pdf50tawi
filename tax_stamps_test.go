package pdf50tawi

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func tinyPNG() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	// fill with red
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			// set pixel
			img.Set(x, y, color.RGBA{R: 255, A: 255})
		}
	}
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

func sampleTaxInfo() TaxInfo {
	return TaxInfo{
		DocumentDetails:      DocumentDetails{BookNumber: "B-001", DocumentNumber: "D-002"},
		Payer:                Payer{TaxID: "1234567890123", TaxID10Digit: "1234567890", Name: "Payer Co.", Address: "123 Main"},
		Payee:                Payee{TaxID: "9876543210987", TaxID10Digit: "0987654321", Name: "John Doe", Address: "99 Road", SequenceNumber: "42", Pnd_1a: true, Pnd_1aSpecial: false, Pnd_2: true, Pnd_2a: false, Pnd_3: true, Pnd_3a: false, Pnd_53: true},
		Income40_1:           IncomeDetail{DatePaid: "01/01/2568", AmountPaid: "100.00", TaxWithheld: "10.00"},
		Income40_2:           IncomeDetail{DatePaid: "01/02/2568", AmountPaid: "200.00", TaxWithheld: "20.00"},
		Income40_3:           IncomeDetail{DatePaid: "01/03/2568", AmountPaid: "300.00", TaxWithheld: "30.00"},
		Income40_4A:          IncomeDetail{DatePaid: "01/04/2568", AmountPaid: "400.00", TaxWithheld: "40.00"},
		Income40_4B_1_1:      IncomeDetail{DatePaid: "01/05/2568", AmountPaid: "500.00", TaxWithheld: "50.00"},
		Income40_4B_1_2:      IncomeDetail{DatePaid: "01/06/2568", AmountPaid: "600.00", TaxWithheld: "60.00"},
		Income40_4B_1_3:      IncomeDetail{DatePaid: "01/07/2568", AmountPaid: "700.00", TaxWithheld: "70.00"},
		Income40_4B_1_4_Rate: "7%",
		Income40_4B_1_4:      IncomeDetail{DatePaid: "01/08/2568", AmountPaid: "800.00", TaxWithheld: "80.00"},
		Income40_4B_2_1:      IncomeDetail{DatePaid: "01/09/2568", AmountPaid: "900.00", TaxWithheld: "90.00"},
		Income40_4B_2_2:      IncomeDetail{DatePaid: "01/10/2568", AmountPaid: "1000.00", TaxWithheld: "100.00"},
		Income40_4B_2_3:      IncomeDetail{DatePaid: "01/11/2568", AmountPaid: "1100.00", TaxWithheld: "110.00"},
		Income40_4B_2_4:      IncomeDetail{DatePaid: "01/12/2568", AmountPaid: "1200.00", TaxWithheld: "120.00"},
		Income40_4B_2_5_Note: "Other note",
		Income40_4B_2_5:      IncomeDetail{DatePaid: "01/13/2568", AmountPaid: "1300.00", TaxWithheld: "130.00"},
		Income5:              IncomeDetail{DatePaid: "01/14/2568", AmountPaid: "1400.00", TaxWithheld: "140.00"},
		Income6:              IncomeDetail{DatePaid: "01/15/2568", AmountPaid: "1500.00", TaxWithheld: "150.00"},
		Income6_Note:         "misc",
		Totals:               Totals{TotalAmountPaid: "4500.00", TotalTaxWithheld: "450.00", TotalTaxWithheldInWords: "four five zero"},
		TotalsInWords:        "four thousand five hundred",
		OtherPayments:        OtherPayments{GovernmentPensionFund: "1", SocialSecurityFund: "2", ProvidentFund: "3"},
		WithholdingType:      WithholdingType{WithholdingTax: true, Forever: false, OneTime: true, Other: true, OtherDetails: "detail"},
		Certification:        Certification{DateOfIssuance: DateOfIssuance{Day: "1", Month: "Jan", Year: "2568"}},
	}
}

func TestTextWatermark(t *testing.T) {
	cfg := TextStmap{Text: "Hello", Dx: 10, Dy: 20, FontSize: 16, FontName: "CustomFont", Position: types.TopLeft}
	wm, err := TextWatermark(cfg)
	if err != nil {
		t.Fatalf("TextWatermark error: %v", err)
	}
	if wm.Dx != cfg.Dx || wm.Dy != cfg.Dy || wm.FontSize != cfg.FontSize || wm.FontName != cfg.FontName || wm.Pos != cfg.Position || !wm.ScaleAbs || wm.OnTop != true {
		t.Fatalf("unexpected watermark fields: %+v", wm)
	}
}

func TestTextWatermark_DefaultFont(t *testing.T) {
	wm, err := TextWatermark(TextStmap{Text: "X", FontSize: 12})
	if err != nil {
		t.Fatal(err)
	}
	if wm.FontName != "THSarabunNew" {
		t.Fatalf("expected default font THSarabunNew, got %s", wm.FontName)
	}
}

func TestImageWatermark(t *testing.T) {
	img := tinyPNG()
	wm, err := ImageWatermark(ImageStamp{Reader: bytes.NewReader(img), Pos: types.BottomRight, Dx: 3, Dy: 4, Scale: 0.5, Opacity: 0.8, OnTop: true})
	if err != nil {
		t.Fatalf("ImageWatermark error: %v", err)
	}
	if wm.Dx != 3 || wm.Dy != 4 || wm.Scale != 0.5 || !wm.ScaleAbs || wm.Opacity != 0.8 || wm.Pos != types.BottomRight || wm.OnTop != true {
		t.Fatalf("unexpected image watermark fields: %+v", wm)
	}
}

func TestPositionHelpers(t *testing.T) {
	st13 := positionTaxID13Digits("1 2 3 4 5 6 7 8 9 0 1 2 3", -10, 14)
	if len(st13) != 13 {
		t.Fatalf("expected 13 stamps, got %d", len(st13))
	}
	st10 := positionTaxID10Digits("0 1 2 3 4 5 6 7 8 9", -20, 12)
	if len(st10) != 10 {
		t.Fatalf("expected 10 stamps, got %d", len(st10))
	}
	st := position("AB", 9, -1, []float64{100})
	if len(st) != 1 || st[0].Text != "A" || st[0].Dx != 100 {
		t.Fatalf("position mismatch: %+v", st)
	}
}

func TestTickAndCheckmarkStamp(t *testing.T) {
	if tick(true) != string(rune(52)) {
		t.Fatalf("tick(true) unexpected")
	}
	if tick(false) != " " {
		t.Fatalf("tick(false) unexpected")
	}
	c := checkmark(true, 1.5, -2.5)
	if c.Text != string(rune(52)) || c.FontName != "ZapfDingbats" || c.FontSize != 10 || c.Dx != 1.5 || c.Dy != -2.5 || c.Position != types.TopLeft {
		t.Fatalf("checkmark stamp mismatch: %+v", c)
	}
}

func TestTextStampsFromTaxInfo(t *testing.T) {
	tax := sampleTaxInfo()
	stamps := TextStampsFromTaxInfo(tax)
	if len(stamps) == 0 {
		t.Fatalf("expected stamps")
	}
	found := map[string]bool{}
	for _, s := range stamps {
		if s.Text == tax.Payer.Name {
			found["payerName"] = true
		}
		if s.Text == tax.Payee.Name {
			found["payeeName"] = true
		}
		if s.Text == tax.DocumentDetails.BookNumber {
			found["book"] = true
		}
		if s.Text == tax.Certification.DateOfIssuance.Day {
			found["day"] = true
		}
	}
	for k, ok := range found {
		if !ok {
			t.Fatalf("missing stamp: %s", k)
		}
	}
}

func TestCertificateImageStamps(t *testing.T) {
	pngBytes := tinyPNG()
	st := CertificateImageStamps(bytes.NewReader(pngBytes), bytes.NewReader(pngBytes))
	if len(st) != 2 {
		t.Fatalf("expected 2 image stamps")
	}
	if st[0].Reader == nil || st[1].Reader == nil {
		t.Fatalf("expected non-nil readers")
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

func TestBuildWriteAndWHTCertificatePDF(t *testing.T) {
	png := tinyPNG()
	// BuildStampedContext
	texts := []TextStmap{{Text: "t", Dx: 10, Dy: -10, FontSize: 12, Position: types.TopLeft}}
	images := []ImageStamp{{Reader: bytes.NewReader(png), Pos: types.BottomLeft, Dx: 5, Dy: 5, Scale: 0.1, Opacity: 1, OnTop: true}}
	ctx, err := BuildStampedContext(mustPDFTemplate(t), texts, images)
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
	if err := WHTCertificatePDF(mustPDFTemplate(t), &out2, sampleTaxInfo(), bytes.NewReader(png), bytes.NewReader(png)); err != nil {
		t.Fatalf("WHTCertificatePDF error: %v", err)
	}
	if out2.Len() == 0 {
		t.Fatalf("expected output bytes for WHTCertificatePDF")
	}
}

func TestIssueWHTCertificatePDF(t *testing.T) {
	png := tinyPNG()
	var out bytes.Buffer
	if err := IssueWHTCertificatePDF(&out, sampleTaxInfo(), bytes.NewReader(png), bytes.NewReader(png)); err != nil {
		t.Fatalf("IssueWHTCertificatePDF error: %v", err)
	}
	if out.Len() == 0 {
		t.Fatalf("expected non-empty output from IssueWHTCertificatePDF")
	}
}
