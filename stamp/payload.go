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

type WithholdingType struct {
	WithholdingTax bool   `json:"withholdingTax"`
	Forever        bool   `json:"forever"`
	OneTime        bool   `json:"oneTime"`
	Other          bool   `json:"other"`
	OtherDetails   string `json:"otherDetails"`
}

type Certification struct {
	PayerSignatureText               string `json:"payerSignatureText"`
	PayerSignatureImageFileLocalPath string `json:"payerSignatureImageFileLocalPath"`
	PayerSignatureImageURL           string `json:"payerSignatureImageURL"`
	PayerSignatureImageBase64        string `json:"payerSignatureImageBase64"`
	JuristicPersonSeal               string `json:"juristicPersonSeal"`
	DateOfIssuance                   string `json:"dateOfIssuance"`
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
			WithholdingTax: true,
			Forever:        true,
			OneTime:        true,
			Other:          true,
			OtherDetails:   "อื่นๆ (อื่นๆ) อื่นๆ อื่นๆ อื่นๆ",
		},
		Certification: Certification{
			PayerSignatureText:               "นายตัวอย่าง ใจดี",
			PayerSignatureImageFileLocalPath: "",
			PayerSignatureImageURL:           "",
			PayerSignatureImageBase64:        "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",
			JuristicPersonSeal:               "[ตราประทับบริษัท]",
			DateOfIssuance:                   "23/08/2568", // format dd/mm/yyyy or any other format just use separator / to separate day, month, year or { day: "", month: "", year: ""}
		},
	}
}
