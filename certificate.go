package pdf50tawi

import (
	"bytes"
	"io"
	"strings"
)

// IssueWHTCertificatePDF generates a filled WHT certificate PDF.
func IssueWHTCertificatePDF(outputPDF io.Writer, taxInfo TaxInfo, sign io.Reader, logo io.Reader) error {
	images := CertificateImageFields(sign, logo)
	texts := TextFieldsFromTaxInfo(taxInfo)
	return fillCertificate(texts, images, outputPDF)
}

// CertificateImageFields returns the positioned image fields for the signature and company seal.
func CertificateImageFields(sign io.Reader, logo io.Reader) []ImageField {
	return []ImageField{
		{Reader: ifNil(sign), Pos: Center, Dx: 86, Dy: -313, Scale: 0.1, Opacity: 1, OnTop: true},
		{Reader: ifNil(logo), Pos: Center, Dx: 212, Dy: -325, Scale: 0.06, Opacity: 1, OnTop: false, Diagonal: 1},
	}
}

// TextFieldsFromTaxInfo converts TaxInfo into the complete set of TextField values
// to be rendered on the certificate form.
func TextFieldsFromTaxInfo(tax TaxInfo) []TextField {

	// Payer Information (ผู้จ่ายเงิน)
	payer := []TextField{
		{Text: tax.Payer.Name, Dx: 58, Dy: -110, FontSize: 14, Position: TopLeft},
		{Text: tax.Payer.Address, Dx: 62, Dy: -132, FontSize: 12, Position: TopLeft},
	}
	payer = append(payer, positionTaxID13Digits(tax.Payer.TaxID, -94, 16)...)
	payer = append(payer, positionTaxID10Digits(tax.Payer.TaxID10Digit, -111, 16)...)

	// Payee Information (ผู้ถูกหักภาษี ณ ที่จ่าย)
	payee := []TextField{
		{Text: tax.Payee.Name, Dx: 58, Dy: -182, FontSize: 14, Position: TopLeft},
		{Text: tax.Payee.Address, Dx: 62, Dy: -208, FontSize: 12, Position: TopLeft},
	}
	payee = append(payee, positionTaxID13Digits(tax.Payee.TaxID, -163, 16)...)
	payee = append(payee, positionTaxID10Digits(tax.Payee.TaxID10Digit, -182, 16)...)
	payee = append(payee, []TextField{
		// Tax Filing Reference (ลำดับที่)
		{Text: tax.Payee.SequenceNumber, Dx: -190, Dy: -236, FontSize: 14, Position: TopCenter},

		checkmark(tax.Payee.Pnd_1a, 211.5, -230),
		checkmark(tax.Payee.Pnd_1aSpecial, 289, -230),
		checkmark(tax.Payee.Pnd_2, 397, -230),
		checkmark(tax.Payee.Pnd_2a, 211.5, -248),
		checkmark(tax.Payee.Pnd_3, 474, -230),
		checkmark(tax.Payee.Pnd_3a, 289, -248),
		checkmark(tax.Payee.Pnd_53, 397, -248),
	}...)

	// Define text fields for Form 50 ทวิ layout
	textFields := []TextField{
		// Document Details (top right)
		{Text: tax.DocumentDetails.BookNumber, Dx: 519, Dy: -59, FontSize: 14, Position: TopLeft},
		{Text: tax.DocumentDetails.DocumentNumber, Dx: 519, Dy: -74, FontSize: 14, Position: TopLeft},

		// Position: Bottom Right
		// Income Details - Row 1 (เงินเดือน ค่าจาง)
		{Text: tax.Income40_1.DatePaid, Dx: 69, Dy: 536, FontSize: 14, Position: BottomCenter},
		{Text: tax.Income40_1.AmountPaid, Dx: -109.5, Dy: 536, FontSize: 14, Position: BottomRight},
		{Text: tax.Income40_1.TaxWithheld, Dx: -38, Dy: 536, FontSize: 14, Position: BottomRight},

		// Income Details - Row 2 (ค่าธรรมเนียม ค่านายหน้า)
		{Text: tax.Income40_2.DatePaid, Dx: 69, Dy: 522, FontSize: 14, Position: BottomCenter},
		{Text: tax.Income40_2.AmountPaid, Dx: -109.5, Dy: 522, FontSize: 14, Position: BottomRight},
		{Text: tax.Income40_2.TaxWithheld, Dx: -38, Dy: 522, FontSize: 14, Position: BottomRight},

		// Income Details - Row 3 (ค่าแห่งลิขสิทธิ์)
		{Text: tax.Income40_3.DatePaid, Dx: 69, Dy: 508, FontSize: 14, Position: BottomCenter},
		{Text: tax.Income40_3.AmountPaid, Dx: -109.5, Dy: 508, FontSize: 14, Position: BottomRight},
		{Text: tax.Income40_3.TaxWithheld, Dx: -38, Dy: 508, FontSize: 14, Position: BottomRight},

		// Income Details - Row 4 (ดอกเบี้ย เงินปันผล) 40 (4) (ก)
		{Text: tax.Income40_4A.DatePaid, Dx: 69, Dy: 494, FontSize: 14, Position: BottomCenter},
		{Text: tax.Income40_4A.AmountPaid, Dx: -109.5, Dy: 494, FontSize: 14, Position: BottomRight},
		{Text: tax.Income40_4A.TaxWithheld, Dx: -38, Dy: 494, FontSize: 14, Position: BottomRight},

		// Income Details - Row 4 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (1) (1.1)
		{Text: tax.Income40_4B_1_1.DatePaid, Dx: 69, Dy: 437, FontSize: 14, Position: BottomCenter},
		{Text: tax.Income40_4B_1_1.AmountPaid, Dx: -109.5, Dy: 437, FontSize: 14, Position: BottomRight},
		{Text: tax.Income40_4B_1_1.TaxWithheld, Dx: -38, Dy: 437, FontSize: 14, Position: BottomRight},

		// Income Details - Row 5 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (1) (1.2)
		{Text: tax.Income40_4B_1_2.DatePaid, Dx: 69, Dy: 420, FontSize: 14, Position: BottomCenter},
		{Text: tax.Income40_4B_1_2.AmountPaid, Dx: -109.5, Dy: 420, FontSize: 14, Position: BottomRight},
		{Text: tax.Income40_4B_1_2.TaxWithheld, Dx: -38, Dy: 420, FontSize: 14, Position: BottomRight},

		// Income Details - Row 5 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (1) (1.3)
		{Text: tax.Income40_4B_1_3.DatePaid, Dx: 69, Dy: 406, FontSize: 14, Position: BottomCenter},
		{Text: tax.Income40_4B_1_3.AmountPaid, Dx: -109.5, Dy: 406, FontSize: 14, Position: BottomRight},
		{Text: tax.Income40_4B_1_3.TaxWithheld, Dx: -38, Dy: 406, FontSize: 14, Position: BottomRight},

		// Income Details - Row 6 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (1) (1.4)
		{Text: tax.Income40_4B_1_4_Rate, Dx: -116, Dy: 390, FontSize: 14, Position: BottomCenter},
		{Text: tax.Income40_4B_1_4.DatePaid, Dx: 69, Dy: 391, FontSize: 14, Position: BottomCenter},
		{Text: tax.Income40_4B_1_4.AmountPaid, Dx: -109.5, Dy: 391, FontSize: 14, Position: BottomRight},
		{Text: tax.Income40_4B_1_4.TaxWithheld, Dx: -38, Dy: 391, FontSize: 14, Position: BottomRight},

		// Income Details - Row 6 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (2) (2.1)
		{Text: tax.Income40_4B_2_1.DatePaid, Dx: 69, Dy: 362, FontSize: 14, Position: BottomCenter},
		{Text: tax.Income40_4B_2_1.AmountPaid, Dx: -109.5, Dy: 362, FontSize: 14, Position: BottomRight},
		{Text: tax.Income40_4B_2_1.TaxWithheld, Dx: -38, Dy: 362, FontSize: 14, Position: BottomRight},

		// Income Details - Row 6 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (2) (2.2)
		{Text: tax.Income40_4B_2_2.DatePaid, Dx: 69, Dy: 333, FontSize: 14, Position: BottomCenter},
		{Text: tax.Income40_4B_2_2.AmountPaid, Dx: -109.5, Dy: 333, FontSize: 14, Position: BottomRight},
		{Text: tax.Income40_4B_2_2.TaxWithheld, Dx: -38, Dy: 333, FontSize: 14, Position: BottomRight},

		// Income Details - Row 6 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (2) (2.3)
		{Text: tax.Income40_4B_2_3.DatePaid, Dx: 69, Dy: 304, FontSize: 14, Position: BottomCenter},
		{Text: tax.Income40_4B_2_3.AmountPaid, Dx: -109.5, Dy: 304, FontSize: 14, Position: BottomRight},
		{Text: tax.Income40_4B_2_3.TaxWithheld, Dx: -38, Dy: 304, FontSize: 14, Position: BottomRight},

		// Income Details - Row 6 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (2) (2.4)
		{Text: tax.Income40_4B_2_4.DatePaid, Dx: 69, Dy: 288, FontSize: 14, Position: BottomCenter},
		{Text: tax.Income40_4B_2_4.AmountPaid, Dx: -109.5, Dy: 288, FontSize: 14, Position: BottomRight},
		{Text: tax.Income40_4B_2_4.TaxWithheld, Dx: -38, Dy: 288, FontSize: 14, Position: BottomRight},

		// Income Details - Row 6 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (2) (2.5)
		{Text: tax.Income40_4B_2_5_Note, Dx: 150, Dy: 275, FontSize: 12, Position: BottomLeft},
		{Text: tax.Income40_4B_2_5.DatePaid, Dx: 69, Dy: 275, FontSize: 14, Position: BottomCenter},
		{Text: tax.Income40_4B_2_5.AmountPaid, Dx: -109.5, Dy: 275, FontSize: 14, Position: BottomRight},
		{Text: tax.Income40_4B_2_5.TaxWithheld, Dx: -38, Dy: 275, FontSize: 14, Position: BottomRight},

		// Income Details - Row 7 5. การจ่ายเงินได้ที่ต้องหักภาษี ณ ที่จ่าย
		{Text: tax.Income5.DatePaid, Dx: 69, Dy: 217, FontSize: 14, Position: BottomCenter},
		{Text: tax.Income5.AmountPaid, Dx: -109.5, Dy: 217, FontSize: 14, Position: BottomRight},
		{Text: tax.Income5.TaxWithheld, Dx: -38, Dy: 217, FontSize: 14, Position: BottomRight},

		// Income Details - Row 8 6. อื่น ๆ (ระบุ)
		{Text: tax.Income6_Note, Dx: 102, Dy: 203, FontSize: 12, Position: BottomLeft},
		{Text: tax.Income6.DatePaid, Dx: 69, Dy: 203, FontSize: 14, Position: BottomCenter},
		{Text: tax.Income6.AmountPaid, Dx: -109.5, Dy: 203, FontSize: 14, Position: BottomRight},
		{Text: tax.Income6.TaxWithheld, Dx: -38, Dy: 203, FontSize: 14, Position: BottomRight},

		// Totals (รวม)
		{Text: tax.Totals.TotalAmountPaid, Dx: -109.5, Dy: 182, FontSize: 14, Position: BottomRight},
		{Text: tax.Totals.TotalTaxWithheld, Dx: -38, Dy: 182, FontSize: 14, Position: BottomRight},
		{Text: tax.Totals.TotalTaxWithheldInWords, Dx: 200, Dy: 163, FontSize: 14, Position: BottomLeft},

		// Other Payments (เงินที่จ่ายเข้ากองทุน)
		{Text: tax.OtherPayments.GovernmentPensionFund, Dx: -318, Dy: 146, FontSize: 12, Position: BottomRight},
		{Text: tax.OtherPayments.SocialSecurityFund, Dx: -190, Dy: 146, FontSize: 12, Position: BottomRight},
		{Text: tax.OtherPayments.ProvidentFund, Dx: -54, Dy: 146, FontSize: 12, Position: BottomRight},

		// Withholding Type
		checkmark(tax.WithholdingType.WithholdingTax, 85, -712),
		checkmark(tax.WithholdingType.Forever, 178, -712),
		checkmark(tax.WithholdingType.OneTime, 285.5, -712),
		checkmark(tax.WithholdingType.Other, 396, -712),
		{Text: tax.WithholdingType.OtherDetails, Dx: 470, Dy: 124, FontSize: 12, Position: BottomLeft},

		// Certification (ลงชื่อ ผู้จ่ายเงิน และวันที่)
		{Text: tax.Certification.DateOfIssuance.Day, Dx: 52, Dy: 77, FontSize: 14, Position: BottomCenter},
		{Text: tax.Certification.DateOfIssuance.Month, Dx: 99, Dy: 77, FontSize: 14, Position: BottomCenter},
		{Text: tax.Certification.DateOfIssuance.Year, Dx: 152, Dy: 77, FontSize: 14, Position: BottomCenter},
	}

	textFields = append(textFields, payer...)
	textFields = append(textFields, payee...)

	return filterEmptyTextFields(textFields)
}

// positionTaxID13Digits creates individual text fields for each digit of a 13-digit tax ID.
func positionTaxID13Digits(taxID string, dy float64, fontSize int) []TextField {
	digits := strings.ReplaceAll(taxID, " ", "")
	xPositions := []float64{378, 396, 408, 420, 432, 450, 463, 474, 486, 498, 517, 529, 548}
	return position(digits, fontSize, dy, xPositions)
}

// positionTaxID10Digits creates individual text fields for each digit of a 10-digit tax ID.
func positionTaxID10Digits(taxID string, dy float64, fontSize int) []TextField {
	digits := strings.ReplaceAll(taxID, " ", "")
	xPositions := []float64{422, 440, 452, 464, 476, 494, 506, 518, 530, 548}
	return position(digits, fontSize, dy, xPositions)
}

func position(digits string, fontSize int, dy float64, xPositions []float64) []TextField {
	var fields []TextField
	for i, digit := range digits {
		if i < len(xPositions) {
			fields = append(fields, TextField{
				Text:     string(digit),
				Dx:       xPositions[i],
				Dy:       dy,
				FontSize: fontSize,
				Position: TopLeft,
			})
		}
	}
	return fields
}

func tick(pnd bool) string {
	if pnd {
		return "✓"
	}
	return ""
}

func checkmark(isSet bool, dx float64, dy float64) TextField {
	return TextField{
		Text:     tick(isSet),
		Dx:       dx,
		Dy:       dy,
		FontSize: 10,
		Position: TopLeft,
	}
}

func filterEmptyTextFields(textFields []TextField) []TextField {
	var filtered []TextField
	for _, field := range textFields {
		if strings.TrimSpace(field.Text) != "" {
			filtered = append(filtered, field)
		}
	}
	return filtered
}

func ifNil(img io.Reader) io.Reader {
	if img == nil {
		return bytes.NewReader(tinyEmptyPNG())
	}
	return img
}
