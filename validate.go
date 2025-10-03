package pdf50tawi

import (
	"fmt"
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
	validateParty(&ve, "payer", t.Payer.Name, t.Payer.TaxID, t.Payer.TaxID10Digit)
	// Payee
	validateParty(&ve, "payee", t.Payee.Name, t.Payee.TaxID, t.Payee.TaxID10Digit)
	// Payee PND validation
	validatePayeePND(&ve, t.Payee)

	if ve.HasErrors() {
		return &ve
	}
	return nil
}

func validatePayeePND(ve *ValidationError, p Payee) {
	if !p.Pnd_1a && !p.Pnd_1aSpecial && !p.Pnd_2 && !p.Pnd_3 && !p.Pnd_2a && !p.Pnd_3a && !p.Pnd_53 {
		ve.Add("ผู้ถูกหักภาษี: ต้องเลือกประเภทเงินได้อย่างน้อยหนึ่งประเภท ภ.ง.ด. 1ก (pnd_1a), ภ.ง.ด. 1ก พิเศษ (pnd_1aSpecial), ภ.ง.ด. 2 (pnd_2), ภ.ง.ด. 3 (pnd_3), ภ.ง.ด. 2ก (pnd_2a), ภ.ง.ด. 3ก (pnd_3a) หรือ ภ.ง.ด. 53 (pnd_53)")
	}
}

func validateParty(ve *ValidationError, prefix, name, tax13, tax10 string) {
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
