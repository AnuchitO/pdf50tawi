package pdf50tawi

import (
	"fmt"
	"net/url"
	"strings"
)

type ValidationError struct {
	Errors []string
}

func (v *ValidationError) Add(msg string)  { v.Errors = append(v.Errors, msg) }
func (v *ValidationError) HasErrors() bool { return len(v.Errors) > 0 }
func (v *ValidationError) Error() string   { return strings.Join(v.Errors, "; ") }

// ValidateTaxInfo validates all fields in TaxInfo and returns a comprehensive error if any.
func ValidateTaxInfo(t TaxInfo) error {
	var ve ValidationError

	// Payer
	ve.validateParty("payer", t.Payer.Name, t.Payer.TaxID, t.Payer.TaxID10Digit)
	// Payee
	ve.validateParty("payee", t.Payee.Name, t.Payee.TaxID, t.Payee.TaxID10Digit)
	// Payee PND validation
	ve.validatePayeePND(t.Payee)
	// Withholding type validation
	ve.validateWithholdingType(t.WithholdingType)

	if ve.HasErrors() {
		return &ve
	}
	return nil
}

func (ve *ValidationError) validatePayeePND(p Payee) {
	if !p.Pnd_1a && !p.Pnd_1aSpecial && !p.Pnd_2 && !p.Pnd_3 && !p.Pnd_2a && !p.Pnd_3a && !p.Pnd_53 {
		ve.Add("ผู้ถูกหักภาษี: ต้องเลือกประเภทเงินได้อย่างน้อยหนึ่งประเภท ภ.ง.ด. 1ก: pnd_1a, ภ.ง.ด. 1ก พิเศษ: pnd_1aSpecial, ภ.ง.ด. 2: pnd_2, ภ.ง.ด. 3: pnd_3, ภ.ง.ด. 2ก: pnd_2a, ภ.ง.ด. 3ก: pnd_3a หรือ ภ.ง.ด. 53: pnd_53")
	}
}

func (ve *ValidationError) validateWithholdingType(w WithholdingType) {
	if !w.WithholdingTax && !w.Forever && !w.OneTime && !w.Other {
		ve.Add("ต้องเลือกประเภทหนังสือรับรองอย่างน้อยหนึ่งประเภท (หัก ณ ที่จ่าย: withholdingTax, ออกให้ตลอดไป: forever, ออกให้ครั้งเดียว: oneTime หรือ อื่น ๆ: other)")
	}
}

func (ve *ValidationError) validateParty(prefix, name, tax13, tax10 string) {
	if strings.TrimSpace(name) == "" {
		ve.Add(fmt.Sprintf("%s.name is required", prefix))
	}
	// Strip spaces first, then check if empty before validating
	strippedTax13 := stripSpaces(tax13)
	strippedTax10 := stripSpaces(tax10)
	if strippedTax13 != "" && !isDigitsLen(strippedTax13, 13) {
		ve.Add(fmt.Sprintf("%s.taxId must be 13 digits", prefix))
	}
	if strippedTax10 != "" && !isDigitsLen(strippedTax10, 10) {
		ve.Add(fmt.Sprintf("%s.taxId10Digit must be 10 digits", prefix))
	}
}

func (ve *ValidationError) validateImage(image Image) {
	// If both source type and value are empty, it's valid (no image)
	if strings.TrimSpace(image.SourceType.String()) == "" && strings.TrimSpace(image.Value) == "" {
		return
	}

	// Validate source type
	switch image.SourceType {
	case SourceTypeUpload:
	case SourceTypeURL:
	case SourceTypeFile:
	default:
		ve.Add("image.sourceType must be 'upload', 'url', or 'file'")
		return
	}

	// Validate value is not empty
	if strings.TrimSpace(image.Value) == "" {
		ve.Add(fmt.Sprintf("image.value is required for '%s' source type", image.SourceType.String()))
		return
	}

	// If URL, validate URL format
	if image.SourceType == SourceTypeURL {
		parsedURL, err := url.Parse(image.Value)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
			ve.Add("image.value must be a valid URL")
		}
	}
}

func stripSpaces(s string) string { return strings.ReplaceAll(s, " ", "") }
func isDigitsLen(s string, l int) bool {
	if len(s) != l {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
