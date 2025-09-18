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
	TaxID          string `json:"taxId"`
	TaxID10Digit   string `json:"taxId10Digit"`
	Name           string `json:"name"`
	Address        string `json:"address"`
	SequenceNumber string `json:"sequenceNumber"`
	Pnd_1a         bool   `json:"pnd_1a"`
	Pnd_1aSpecial  bool   `json:"pnd_1aSpecial"`
	Pnd_2          bool   `json:"pnd_2"`
	Pnd_3          bool   `json:"pnd_3"`
	Pnd_2a         bool   `json:"pnd_2a"`
	Pnd_3a         bool   `json:"pnd_3a"`
	Pnd_53         bool   `json:"pnd_53"`
}

type IncomeDetail struct {
	DatePaid    string `json:"datePaid"`
	AmountPaid  string `json:"amountPaid"`
	TaxWithheld string `json:"taxWithheld"`
}

type IncomeDetails struct {
	// 1. เงินเดือน ค่าจาง เบี้ยเลี้ยง โบนัส ฯลฯ ตามมาตรา 40 (1)
	Income40_1 IncomeDetail `json:"income40_1"`
	// 2. ค่าธรรมเนียม ค่านายหน้า ฯลฯ ตามมาตรา 40 (2)
	Income40_2 IncomeDetail `json:"income40_2"`
	// 3. ค่าแห่งลิขสิทธิ์ ฯลฯ ตามมาตรา 40 (3)
	Income40_3 IncomeDetail `json:"income40_3"`
	// 4. (ก) ดอกเบี้ย ฯลฯ ตามมาตรา 40 (4) (ก)
	Income40_4A IncomeDetail `json:"income40_4a"`
	// 4. (ข) เงินปันผล เงินส่วนแบ่งกำไร ฯลฯ ตามมาตรา 40 (4) (ข)
	// (1) (1.1)
	Income40_4B_1_1 IncomeDetail `json:"income40_4b_1_1"`
	// (1) (1.2)
	Income40_4B_1_2 IncomeDetail `json:"income40_4b_1_2"`
	// (1) (1.3)
	Income40_4B_1_3 IncomeDetail `json:"income40_4b_1_3"`
	// (2) (2.1)
	Income40_4B_2_1 IncomeDetail `json:"income40_4b_2_1"`
	// (2) (2.2)
	Income40_4B_2_2 IncomeDetail `json:"income40_4b_2_2"`
	// (2) (2.3)
	Income40_4B_2_3 IncomeDetail `json:"income40_4b_2_3"`
	// (2) (2.4)
	Income40_4B_2_4 IncomeDetail `json:"income40_4b_2_4"`
	// (2) (2.5)
	Income40_4B_2_5 IncomeDetail `json:"income40_4b_2_5"`
	// 5. การจ่ายเงินได้ที่ต้องหักภาษี ณ ที่จ่าย
	Income40_5 IncomeDetail `json:"income40_5"`
	// 6. อื่น ๆ (ระบุ)
	Income40_6 IncomeDetail `json:"income40_6"`
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
	DocumentDetails DocumentDetails `json:"documentDetails"`
	Payer           Payer           `json:"payer"`
	Payee           Payee           `json:"payee"`
	Income40_1      IncomeDetail    `json:"income40_1"`
	Income40_2      IncomeDetail    `json:"income40_2"`
	Income40_3      IncomeDetail    `json:"income40_3"`
	Income40_4A     IncomeDetail    `json:"income40_4A"`
	Income40_4B_1_1 IncomeDetail    `json:"income40_4B_1_1"`
	Income40_4B_1_2 IncomeDetail    `json:"income40_4B_1_2"`
	Income40_4B_1_3 IncomeDetail    `json:"income40_4B_1_3"`
	Income40_4B_1_4 IncomeDetail    `json:"income40_4B_1_4"`
	Income40_4B_2_1 IncomeDetail    `json:"income40_4B_2_1"`
	Income40_4B_2_2 IncomeDetail    `json:"income40_4B_2_2"`
	Income40_4B_2_3 IncomeDetail    `json:"income40_4B_2_3"`
	Income40_4B_2_4 IncomeDetail    `json:"income40_4B_2_4"`
	Income40_4B_2_5 IncomeDetail    `json:"income40_4B_2_5"`
	Income40_5      IncomeDetail    `json:"income40_5"`
	Income40_6      IncomeDetail    `json:"income40_6"`
	Totals          Totals          `json:"totals"`
	OtherPayments   OtherPayments   `json:"otherPayments"`
	WithholdingType WithholdingType `json:"withholdingType"`
	Certification   Certification   `json:"certification"`
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
			TaxID:          "3210987654321",
			TaxID10Digit:   "1234567890",
			Name:           "นางสาวสมชาย นามสกุลยาวมากไหมนะก็ไม่รู้เหมือนกัน",
			Address:        "555 ต.ทุ่งนา  อ.ทุ่งนา  จ.ชลบุรี  12345",
			SequenceNumber: "001",
			Pnd_1a:         true,
			Pnd_1aSpecial:  true,
			Pnd_2:          true,
			Pnd_3:          true,
			Pnd_2a:         true,
			Pnd_3a:         true,
			Pnd_53:         true,
		},
		Income40_1:      IncomeDetail{DatePaid: "01/01/2568", AmountPaid: "401,010.01", TaxWithheld: "12,030.30"},
		Income40_2:      IncomeDetail{DatePaid: "02/02/2568", AmountPaid: "402,020.02", TaxWithheld: "12,060.60"},
		Income40_3:      IncomeDetail{DatePaid: "03/03/2568", AmountPaid: "403,030.03", TaxWithheld: "12,090.90"},
		Income40_4A:     IncomeDetail{DatePaid: "04/04/2568", AmountPaid: "404,040.04", TaxWithheld: "12,121.20"},
		Income40_4B_1_1: IncomeDetail{DatePaid: "11/04/2568", AmountPaid: "404,101.05", TaxWithheld: "12,123.03"},
		Income40_4B_1_2: IncomeDetail{DatePaid: "12/04/2568", AmountPaid: "404,102.06", TaxWithheld: "12,123.06"},
		Income40_4B_1_3: IncomeDetail{DatePaid: "13/04/2568", AmountPaid: "404,103.07", TaxWithheld: "12,123.09"},
		Income40_4B_1_4: IncomeDetail{DatePaid: "14/04/2568", AmountPaid: "404,104.08", TaxWithheld: "12,123.12"},

		Income40_4B_2_1: IncomeDetail{DatePaid: "21/04/2568", AmountPaid: "404,201.09", TaxWithheld: "12,126.03"},
		Income40_4B_2_2: IncomeDetail{DatePaid: "22/04/2568", AmountPaid: "404,202.10", TaxWithheld: "12,126.06"},
		Income40_4B_2_3: IncomeDetail{DatePaid: "23/04/2568", AmountPaid: "404,203.11", TaxWithheld: "12,126.09"},
		Income40_4B_2_4: IncomeDetail{DatePaid: "24/04/2568", AmountPaid: "404,204.12", TaxWithheld: "12,126.12"},
		Income40_4B_2_5: IncomeDetail{DatePaid: "25/04/2568", AmountPaid: "404,205.13", TaxWithheld: "12,126.15"},

		Income40_5: IncomeDetail{DatePaid: "05/05/2568", AmountPaid: "500,555.55", TaxWithheld: "15,555.55"},
		Income40_6: IncomeDetail{DatePaid: "06/06/2568", AmountPaid: "600,666.66", TaxWithheld: "16,666.66"},
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
