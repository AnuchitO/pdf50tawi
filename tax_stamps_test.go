package pdf50tawi

import (
	"bytes"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

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
		Income40_4B_1_4_Rate: "ร้อยละ 7",
		Income40_4B_1_4:      IncomeDetail{DatePaid: "01/08/2568", AmountPaid: "800.00", TaxWithheld: "80.00"},
		Income40_4B_2_1:      IncomeDetail{DatePaid: "01/09/2568", AmountPaid: "900.00", TaxWithheld: "90.00"},
		Income40_4B_2_2:      IncomeDetail{DatePaid: "01/10/2568", AmountPaid: "1000.00", TaxWithheld: "100.00"},
		Income40_4B_2_3:      IncomeDetail{DatePaid: "01/11/2568", AmountPaid: "1100.00", TaxWithheld: "110.00"},
		Income40_4B_2_4:      IncomeDetail{DatePaid: "01/12/2568", AmountPaid: "1200.00", TaxWithheld: "120.00"},
		Income40_4B_2_5_Note: "ใส่หมายเหตุ",
		Income40_4B_2_5:      IncomeDetail{DatePaid: "01/13/2568", AmountPaid: "1300.00", TaxWithheld: "130.00"},
		Income5:              IncomeDetail{DatePaid: "01/14/2568", AmountPaid: "1400.00", TaxWithheld: "140.00"},
		Income6:              IncomeDetail{DatePaid: "01/15/2568", AmountPaid: "1500.00", TaxWithheld: "150.00"},
		Income6_Note:         "ใส่หมายเหตุ",
		Totals:               Totals{TotalAmountPaid: "4500.00", TotalTaxWithheld: "450.00", TotalTaxWithheldInWords: "สี่ร้อยห้าสิบบาทถ้วน"},
		OtherPayments:        OtherPayments{GovernmentPensionFund: "1", SocialSecurityFund: "2", ProvidentFund: "3"},
		WithholdingType:      WithholdingType{WithholdingTax: true, Forever: false, OneTime: true, Other: true, OtherDetails: "detail"},
		Certification:        Certification{DateOfIssuance: DateOfIssuance{Day: "1", Month: "Jan", Year: "2568"}},
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
	t.Run("tick", func(t *testing.T) {
		if tick(true) != string(rune(52)) {
			t.Fatalf("tick(true) should return tick")
		}
		if tick(false) != " " {
			t.Fatalf("tick(false) should return space")
		}
	})

	t.Run("checkmark", func(t *testing.T) {
		c := checkmark(true, 1.5, -2.5)
		if c.Text != string(rune(52)) || c.FontName != "ZapfDingbats" || c.FontSize != 10 || c.Dx != 1.5 || c.Dy != -2.5 || c.Position != types.TopLeft {
			t.Fatalf("checkmark stamp mismatch: %+v", c)
		}
	})

	t.Run("checkmark with default values", func(t *testing.T) {
		c := checkmark(false, 0, 0)
		if c.Text != " " || c.FontName != "ZapfDingbats" || c.FontSize != 10 || c.Dx != 0 || c.Dy != 0 || c.Position != types.TopLeft {
			t.Fatalf("checkmark stamp mismatch: %+v", c)
		}
	})
}

func TestTextStampsFromTaxInfo(t *testing.T) {
	tax := sampleTaxInfo()
	stamps := TextStampsFromTaxInfo(tax)

	// Check exact count
	if len(stamps) != 122 {
		t.Fatalf("expected 122 fields, got %d fields: that is all fields in form 50 tawi", len(stamps))
	}
}

func TestCertificateImageStamps(t *testing.T) {
	pngBytes := tinyEmptyPNG()
	st := CertificateImageStamps(bytes.NewReader(pngBytes), bytes.NewReader(pngBytes))
	if len(st) != 2 {
		t.Fatalf("expected 2 image stamps")
	}
	if st[0].Reader == nil || st[1].Reader == nil {
		t.Fatalf("expected non-nil readers")
	}
}

func TestCertificateImageStampsWithNilInputs(t *testing.T) {
	// Test the ifNil fallback case - when inputs are nil
	st := CertificateImageStamps(nil, nil)
	if len(st) != 2 {
		t.Fatalf("expected 2 image stamps even with nil inputs")
	}

	// Verify that nil inputs are replaced with valid readers
	if st[0].Reader == nil {
		t.Fatalf("expected signature reader to be non-nil (should fallback to empty PNG)")
	}
	if st[1].Reader == nil {
		t.Fatalf("expected logo reader to be non-nil (should fallback to empty PNG)")
	}

	// Verify the readers contain valid PNG data by checking length
	signBytes := make([]byte, 100)
	n, _ := st[0].Reader.Read(signBytes)
	if n < 10 { // Should have at least PNG signature
		t.Fatalf("expected valid PNG data in signature reader, got %d bytes", n)
	}

	logoBytes := make([]byte, 100)
	n, _ = st[1].Reader.Read(logoBytes)
	if n < 10 { // Should have at least PNG signature
		t.Fatalf("expected valid PNG data in logo reader, got %d bytes", n)
	}
}
