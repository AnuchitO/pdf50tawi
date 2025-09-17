package main

type DocumentDetails struct {
	BookNumber     string `json:"bookNumber"`
	DocumentNumber string `json:"documentNumber"`
}

type Payer struct {
	TaxID        string `json:"taxId"`
	TaxID10Digit string `json:"taxId10Digit"`
	Name         string `json:"name"`
	Address      string `json:"address"`
}

type Payee struct {
	TaxID        string `json:"taxId"`
	TaxID10Digit string `json:"taxId10Digit"`
	Name         string `json:"name"`
	Address      string `json:"address"`
}

type Form struct {
	Selected bool `json:"selected"`
}

type TaxFilingReference struct {
	SequenceNumber string `json:"sequenceNumber"`
	Forms          []Form `json:"forms"`
}

type IncomeDetail struct {
	DatePaid    string `json:"datePaid"`
	AmountPaid  string `json:"amountPaid"`
	TaxWithheld string `json:"taxWithheld"`
}

type Totals struct {
	TotalAmountPaid         string `json:"totalAmountPaid"`
	TotalTaxWithheld        string `json:"totalTaxWithheld"`
	TotalTaxWithheldInWords string `json:"totalTaxWithheldInWords"`
}

type OtherPayments struct {
	GovernmentPensionFund string `json:"governmentPensionFund"`
	SocialSecurityFund    string `json:"socialSecurityFund"`
	ProvidentFund         string `json:"providentFund"`
}

type WithholdingOption struct {
	Selected bool   `json:"selected"`
	Details  string `json:"details,omitempty"`
}

type WithholdingType struct {
	Options []WithholdingOption `json:"options"`
}

type Certification struct {
	PayerSignature     string `json:"payerSignature"`
	JuristicPersonSeal string `json:"juristicPersonSeal"`
	DateOfIssuance     string `json:"dateOfIssuance"`
}

type Payload struct {
	DocumentDetails    DocumentDetails    `json:"documentDetails"`
	Payer              Payer              `json:"payer"`
	Payee              Payee              `json:"payee"`
	TaxFilingReference TaxFilingReference `json:"taxFilingReference"`
	IncomeDetails      []IncomeDetail     `json:"incomeDetails"`
	Totals             Totals             `json:"totals"`
	OtherPayments      OtherPayments      `json:"otherPayments"`
	WithholdingType    WithholdingType    `json:"withholdingType"`
	Certification      Certification      `json:"certification"`
}

// demo payload
func DemoPayload() Payload {
	return Payload{
		DocumentDetails: DocumentDetails{
			BookNumber:     "001",
			DocumentNumber: "2568-001",
		},
		Payer: Payer{
			TaxID:        "1234567890123",
			TaxID10Digit: "0987654321",
			Name:         "บริษัท ตัวอย่าง จำกัด",
			Address:      "123 ถนนสุขุมวิท แขวงคลองตัน เขตวัฒนา กรุงเทพฯ 10110",
		},
		Payee: Payee{
			TaxID:        "3210987654321",
			TaxID10Digit: "1234567890",
			Name:         "นางสาวสมชาย นามสกุลยาวมากไหมนะก็ไม่รู้เหมือนกัน",
			Address:      "555 ต.ทุ่งนา  อ.ทุ่งนา  จ.ชลบุรี  12345",
		},
		TaxFilingReference: TaxFilingReference{
			SequenceNumber: "001",
			Forms: []Form{
				{Selected: true},
			},
		},
		IncomeDetails: []IncomeDetail{
			{DatePaid: "01/01/2568", AmountPaid: "401,010.01", TaxWithheld: "12,030.30"},
			{DatePaid: "02/02/2568", AmountPaid: "402,020.02", TaxWithheld: "12,060.60"},
			{DatePaid: "03/03/2568", AmountPaid: "403,030.03", TaxWithheld: "12,090.90"},
			{DatePaid: "04/04/2568", AmountPaid: "404,040.04", TaxWithheld: "12,121.20"},
			{DatePaid: "11/04/2568", AmountPaid: "404,101.05", TaxWithheld: "12,123.03"},
			{DatePaid: "12/04/2568", AmountPaid: "404,102.06", TaxWithheld: "12,123.06"},
			{DatePaid: "13/04/2568", AmountPaid: "404,103.07", TaxWithheld: "12,123.09"},
			{DatePaid: "14/04/2568", AmountPaid: "404,104.08", TaxWithheld: "12,123.12"},
			{DatePaid: "21/04/2568", AmountPaid: "404,201.09", TaxWithheld: "12,126.03"},
			{DatePaid: "22/04/2568", AmountPaid: "404,202.10", TaxWithheld: "12,126.06"},
			{DatePaid: "23/04/2568", AmountPaid: "404,203.11", TaxWithheld: "12,126.09"},
			{DatePaid: "24/04/2568", AmountPaid: "404,204.12", TaxWithheld: "12,126.12"},
			{DatePaid: "25/04/2568", AmountPaid: "404,205.13", TaxWithheld: "12,126.15"},
			{DatePaid: "05/05/2568", AmountPaid: "500,555.14", TaxWithheld: "15,016.65"},
			{DatePaid: "06/06/2568", AmountPaid: "600,666.15", TaxWithheld: "18,019.98"},
		},
		Totals: Totals{
			TotalAmountPaid:         "5,847,066.90",
			TotalTaxWithheld:        "175,411.22",
			TotalTaxWithheldInWords: "หนึ่งแสนเจ็ดหมื่นห้าพันสี่ร้อยสิบเอ็ดบาทยี่สิบสองสตางค์",
		},
		OtherPayments: OtherPayments{
			GovernmentPensionFund: "22,222.22",
			SocialSecurityFund:    "33,333.33",
			ProvidentFund:         "44,444.44",
		},
		WithholdingType: WithholdingType{
			Options: []WithholdingOption{
				{Selected: true, Details: "ออกให้ตลอดไป"},
				{Selected: false, Details: "ออกให้ครั้งเดียว"},
				{Selected: false, Details: "อื่นๆ (ระบุ)"},
			},
		},
		Certification: Certification{
			PayerSignature:     "นายตัวอย่าง ใจดี",
			JuristicPersonSeal: "[ตราประทับบริษัท]",
			DateOfIssuance:     "99/09/2568",
		},
	}
}
