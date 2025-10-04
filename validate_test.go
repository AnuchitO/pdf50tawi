package pdf50tawi

import (
	"strings"
	"testing"
)

func validTaxInfo() TaxInfo {
	return TaxInfo{
		DocumentDetails: DocumentDetails{BookNumber: "001", DocumentNumber: "WHT-001012568"},
		Payer:           Payer{Name: "บริษัท ตัวอย่าง จำกัด", Address: "123 ถนนสุขุมวิท แขวงคลองตัน เขตวัฒนา กรุงเทพฯ 10110", TaxID: "1234567890123", TaxID10Digit: "1234567890"},
		Payee:           Payee{Name: "นางสาวสมชาย นามสกุลยาวมากไหมนะก็ไม่รู้เหมือนกัน", Address: "555 ต.ทุ่งนา  อ.ทุ่งนา  จ.ชลบุรี  12345", TaxID: "9876543210987", TaxID10Digit: "0987654321", SequenceNumber: "1", Pnd_1a: true},
		Income40_1:      IncomeDetail{DatePaid: "01 มกราคม 2568", AmountPaid: "1000.00", TaxWithheld: "30.00"},
		Totals:          Totals{TotalAmountPaid: "1000.00", TotalTaxWithheld: "30.00", TotalTaxWithheldInWords: "สามสิบบาทถ้วน"},
		WithholdingType: WithholdingType{WithholdingTax: true},
		Certification:   Certification{DateOfIssuance: DateOfIssuance{Day: "1", Month: "มกราคม", Year: "2568"}},
	}
}

func TestValidateTaxInfo_OK(t *testing.T) {
	v := validTaxInfo()
	if err := ValidateTaxInfo(v); err != nil {
		t.Fatalf("expected valid tax info, got error: %v", err)
	}
}

func TestValidateTaxID(t *testing.T) {
	testCases := []struct {
		name             string
		input            string
		expectedErrorMsg string
	}{
		// Success cases
		{"Valid13Digits", "1234567890123", ""},
		{"ValidWithSpaces", "123 456 789 01 23", ""},
		{"ValidWithMultipleSpaces", "12 345 678 90 123", ""},

		// Failure cases
		{"TooShort_3Digits", "123", "payer.taxId must be 13 digits"},
		{"TooShort_12Digits", "123456789012", "payer.taxId must be 13 digits"},
		{"TooLong_14Digits", "12345678901234", "payer.taxId must be 13 digits"},
		{"TooLong_15Digits", "123456789012345", "payer.taxId must be 13 digits"},

		// Edge cases
		{"EmptyString", "", ""},
		{"Only13Spaces", "             ", ""},
		{"NonDigits", "12345678901a3", "payer.taxId must be 13 digits"},
		{"NonDigitsWithSpaces", "123 456 789 a1 23", "payer.taxId must be 13 digits"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := validTaxInfo()
			v.Payer.TaxID = tc.input
			err := ValidateTaxInfo(v)

			if tc.expectedErrorMsg == "" {
				if err != nil {
					t.Fatalf("expected no error for input %q, got %v", tc.input, err)
				}
			} else {
				if err == nil {
					t.Fatalf("expected error for input %q, got nil", tc.input)
				}
				if !strings.Contains(err.Error(), tc.expectedErrorMsg) {
					t.Fatalf("expected error to contain %q, got %q", tc.expectedErrorMsg, err.Error())
				}
			}
		})
	}
}

func TestValidateTaxID10Digit(t *testing.T) {
	testCases := []struct {
		name             string
		input            string
		expectedErrorMsg string
	}{
		// Success cases
		{"Valid10Digits", "1234567890", ""},
		{"ValidWithSpaces", "12 345 678 90", ""},
		{"ValidWithMultipleSpaces", "1 234 567 890", ""},

		// Failure cases
		{"TooShort_9Digits", "123456789", "payer.taxId10Digit must be 10 digits"},
		{"TooShort_5Digits", "12345", "payer.taxId10Digit must be 10 digits"},
		{"TooLong_11Digits", "12345678901", "payer.taxId10Digit must be 10 digits"},
		{"TooLong_12Digits", "123456789012", "payer.taxId10Digit must be 10 digits"},

		// Edge cases
		{"EmptyString", "", ""},
		{"Only10Spaces", "          ", ""},
		{"NonDigits", "123456789a", "payer.taxId10Digit must be 10 digits"},
		{"NonDigitsWithSpaces", "123 456 78 9a", "payer.taxId10Digit must be 10 digits"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := validTaxInfo()
			v.Payer.TaxID10Digit = tc.input
			err := ValidateTaxInfo(v)

			if tc.expectedErrorMsg == "" {
				if err != nil {
					t.Fatalf("expected no error for input %q, got %v", tc.input, err)
				}
			} else {
				if err == nil {
					t.Fatalf("expected error for input %q, got nil", tc.input)
				}
				if !strings.Contains(err.Error(), tc.expectedErrorMsg) {
					t.Fatalf("expected error to contain %q, got %q", tc.expectedErrorMsg, err.Error())
				}
			}
		})
	}
}

func TestValidatePayeePND(t *testing.T) {
	var ve ValidationError

	testCases := []struct {
		name             string
		payee            Payee
		expectedErrorMsg string
	}{
		// Success cases - at least one PND field is true
		{"Valid_ภ.ง.ด. 1ก", Payee{Pnd_1a: true}, ""},
		{"Valid_ภ.ง.ด. 1ก พิเศษ", Payee{Pnd_1aSpecial: true}, ""},
		{"Valid_ภ.ง.ด. 2", Payee{Pnd_2: true}, ""},
		{"Valid_ภ.ง.ด. 3", Payee{Pnd_3: true}, ""},
		{"Valid_ภ.ง.ด. 2ก", Payee{Pnd_2a: true}, ""},
		{"Valid_ภ.ง.ด. 3ก", Payee{Pnd_3a: true}, ""},
		{"Valid_ภ.ง.ด. 53", Payee{Pnd_53: true}, ""},

		// Failure case - no PND fields set
		{"Invalid_NoPND", Payee{}, "ผู้ถูกหักภาษี: ต้องเลือกประเภทเงินได้อย่างน้อยหนึ่งประเภท"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset validation error for each test
			ve = ValidationError{}

			// Test validatePayeePND directly
			ve.validatePayeePND(tc.payee)

			if tc.expectedErrorMsg == "" {
				// Expect no error
				if ve.HasErrors() {
					t.Fatalf("expected no error for %s, got %v", tc.name, ve.Error())
				}
			} else {
				// Expect specific error
				if !ve.HasErrors() {
					t.Fatalf("expected error for %s, got no error", tc.name)
				}
				if !strings.Contains(ve.Error(), tc.expectedErrorMsg) {
					t.Fatalf("expected error to contain %q, got %q", tc.expectedErrorMsg, ve.Error())
				}
			}
		})
	}
}

func TestValidateWithholdingType(t *testing.T) {
	var ve ValidationError

	testCases := []struct {
		name             string
		withholdingType  WithholdingType
		expectedErrorMsg string
	}{
		// Success cases - at least one field is true
		{"Valid_WithholdingTax", WithholdingType{WithholdingTax: true}, ""},
		{"Valid_Forever", WithholdingType{Forever: true}, ""},
		{"Valid_OneTime", WithholdingType{OneTime: true}, ""},
		{"Valid_Other", WithholdingType{Other: true}, ""},
		{"Valid_Multiple", WithholdingType{WithholdingTax: true, Forever: true}, ""},

		// Failure case - no fields set
		{"Invalid_NoneSelected", WithholdingType{}, "ต้องเลือกประเภทหนังสือรับรองอย่างน้อยหนึ่งประเภท"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset validation error for each test
			ve = ValidationError{}

			// Test validateWithholdingType directly
			ve.validateWithholdingType(tc.withholdingType)

			if tc.expectedErrorMsg == "" {
				// Expect no error
				if ve.HasErrors() {
					t.Fatalf("expected no error for %s, got %v", tc.name, ve.Error())
				}
			} else {
				// Expect specific error
				if !ve.HasErrors() {
					t.Fatalf("expected error for %s, got no error", tc.name)
				}
				if !strings.Contains(ve.Error(), tc.expectedErrorMsg) {
					t.Fatalf("expected error to contain %q, got %q", tc.expectedErrorMsg, ve.Error())
				}
			}
		})
	}
}

func TestValidateTaxInfo_AllErrors(t *testing.T) {
	v := TaxInfo{
		Payer: Payer{
			TaxID:        "000",
			TaxID10Digit: "000",
			Name:         "",
		},
		Payee: Payee{
			TaxID:        "000",
			TaxID10Digit: "000",
			Name:         "",
		},
		WithholdingType: WithholdingType{},
	}

	err := ValidateTaxInfo(v)
	if err == nil {
		t.Fatalf("expected multiple validation errors, got nil")
	}

	msg := err.Error()
	expectedErrors := []string{
		"payer.name is required",
		"payee.name is required",
		"payee.taxId must be 13 digits",
		"payee.taxId10Digit must be 10 digits",
		"ผู้ถูกหักภาษี: ต้องเลือกประเภทเงินได้อย่างน้อยหนึ่งประเภท",
		"ต้องเลือกประเภทหนังสือรับรองอย่างน้อยหนึ่งประเภท",
	}

	for _, expected := range expectedErrors {
		if !strings.Contains(msg, expected) {
			t.Fatalf("expected error to contain %q, got %q", expected, msg)
		}
	}
}

func TestValidateImage(t *testing.T) {
	var ve ValidationError

	// sourctype will be empty or url or file or upload
	testCases := []struct {
		name             string
		image            Image
		expectedErrorMsg string
	}{
		// Success cases
		{"Valid_Url", Image{SourceType: SourceTypeURL, Value: "https://example.com/image.png"}, ""},
		{"Valid_File", Image{SourceType: SourceTypeFile, Value: "path/to/image.png"}, ""},
		{"Valid_Upload", Image{SourceType: SourceTypeUpload, Value: "path/to/image.png"}, ""},
		{"Valid_SourceType_empty_value_should_be_empty", Image{SourceType: "", Value: ""}, ""},
		// // Failure cases
		{"Invalid_SourceType", Image{SourceType: "invalid", Value: "path/to/image.png"}, "image.sourceType must be 'upload', 'url', or 'file'"},
		{"Invalid_Value", Image{SourceType: SourceTypeURL, Value: ""}, "image.value is required for 'url' source type"},
		{"Invalid_Url", Image{SourceType: SourceTypeURL, Value: "invalid-url"}, "image.value must be a valid URL"},
		{"Invalid_SourceType_empty_value_not_be_empty", Image{SourceType: "", Value: "path/to/image.png"}, "image.sourceType must be 'upload', 'url', or 'file'"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset validation error for each test
			ve = ValidationError{}

			// Test validateImage directly
			ve.validateImage(tc.image)

			if tc.expectedErrorMsg == "" {
				// Expect no error
				if ve.HasErrors() {
					t.Fatalf("expected no error for %s, got %v", tc.name, ve.Error())
				}
			} else {
				// Expect specific error
				if !ve.HasErrors() {
					t.Fatalf("expected error for %s, got no error", tc.name)
				}
				if !strings.Contains(ve.Error(), tc.expectedErrorMsg) {
					t.Fatalf("expected error to contain %q, got %q", tc.expectedErrorMsg, ve.Error())
				}
			}
		})
	}
}
