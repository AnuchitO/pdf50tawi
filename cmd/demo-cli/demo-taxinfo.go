package main

import "github.com/anuchito/pdf50tawi"

// demo tax info
func demoTaxInfo() pdf50tawi.TaxInfo {
	return pdf50tawi.TaxInfo{
		DocumentDetails: pdf50tawi.DocumentDetails{
			BookNumber:     "001",
			DocumentNumber: "2568-001",
		},
		Payer: pdf50tawi.Payer{
			TaxID:        "1234567890123",
			TaxID10Digit: "0987654321",
			Name:         "บริษัท ตัวอย่าง จำกัด",
			Address:      "123 ถนนสุขุมวิท แขวงคลองตัน เขตวัฒนา กรุงเทพฯ 10110",
		},
		Payee: pdf50tawi.Payee{
			TaxID:          "3210987654321",
			TaxID10Digit:   "1234567890",
			Name:           "นางสาวสมชาย นามสกุลยาวมากไหมนะก็ไม่รู้เหมือนกัน",
			Address:        "555 ต.ทุ่งนา  อ.ทุ่งนา  จ.ชลบุรี  12345",
			SequenceNumber: "321",
			Pnd_1a:         true,
			Pnd_1aSpecial:  true,
			Pnd_2:          true,
			Pnd_3:          true,
			Pnd_2a:         true,
			Pnd_3a:         true,
			Pnd_53:         true,
		},
		Income40_1:           pdf50tawi.IncomeDetail{DatePaid: "01 มกราคม 2568", AmountPaid: "401,010.01", TaxWithheld: "12,030.30"},
		Income40_2:           pdf50tawi.IncomeDetail{DatePaid: "02 ก.พ. 2568", AmountPaid: "402,020.02", TaxWithheld: "12,060.60"},
		Income40_3:           pdf50tawi.IncomeDetail{DatePaid: "03/03/2568", AmountPaid: "403,030.03", TaxWithheld: "12,090.90"},
		Income40_4A:          pdf50tawi.IncomeDetail{DatePaid: "04/04/2568", AmountPaid: "404,040.04", TaxWithheld: "12,121.20"},
		Income40_4B_1_1:      pdf50tawi.IncomeDetail{DatePaid: "11/04/2568", AmountPaid: "404,101.05", TaxWithheld: "12,123.03"},
		Income40_4B_1_2:      pdf50tawi.IncomeDetail{DatePaid: "12/04/2568", AmountPaid: "404,102.06", TaxWithheld: "12,123.06"},
		Income40_4B_1_3:      pdf50tawi.IncomeDetail{DatePaid: "13/04/2568", AmountPaid: "404,103.07", TaxWithheld: "12,123.09"},
		Income40_4B_1_4_Rate: "ร้อยละ 10",
		Income40_4B_1_4:      pdf50tawi.IncomeDetail{DatePaid: "14/04/2568", AmountPaid: "404,104.08", TaxWithheld: "12,123.12"},

		Income40_4B_2_1:      pdf50tawi.IncomeDetail{DatePaid: "21/04/2568", AmountPaid: "404,201.09", TaxWithheld: "12,126.03"},
		Income40_4B_2_2:      pdf50tawi.IncomeDetail{DatePaid: "22/04/2568", AmountPaid: "404,202.10", TaxWithheld: "12,126.06"},
		Income40_4B_2_3:      pdf50tawi.IncomeDetail{DatePaid: "23/04/2568", AmountPaid: "404,203.11", TaxWithheld: "12,126.09"},
		Income40_4B_2_4:      pdf50tawi.IncomeDetail{DatePaid: "24/04/2568", AmountPaid: "404,204.12", TaxWithheld: "12,126.12"},
		Income40_4B_2_5_Note: "ระบุ ใดๆๆๆๆๆ ",
		Income40_4B_2_5:      pdf50tawi.IncomeDetail{DatePaid: "25/04/2568", AmountPaid: "404,205.13", TaxWithheld: "12,126.15"},

		Income5:      pdf50tawi.IncomeDetail{DatePaid: "05/05/2568", AmountPaid: "500,555.55", TaxWithheld: "15,555.55"},
		Income6_Note: "อื่นๆ (อื่นๆ) อื่นๆ อื่นๆ อื่นๆ",
		Income6:      pdf50tawi.IncomeDetail{DatePaid: "06/06/2568", AmountPaid: "600,666.66", TaxWithheld: "16,666.66"},
		Totals: pdf50tawi.Totals{
			TotalAmountPaid:         "5,847,066.90",
			TotalTaxWithheld:        "175,411.22",
			TotalTaxWithheldInWords: "หนึ่งแสนเจ็ดหมื่นห้าพันสี่ร้อยสิบเอ็ดบาทยี่สิบสองสตางค์",
		},
		OtherPayments: pdf50tawi.OtherPayments{
			GovernmentPensionFund: "22,222.22",
			SocialSecurityFund:    "33,333.33",
			ProvidentFund:         "44,444.44",
		},
		WithholdingType: pdf50tawi.WithholdingType{
			WithholdingTax: true,
			Forever:        true,
			OneTime:        true,
			Other:          true,
			OtherDetails:   "อื่นๆ (อื่นๆ) อื่นๆ อื่นๆ อื่นๆ",
		},
		Certification: pdf50tawi.Certification{
			PayerSignatureImage: pdf50tawi.Image{
				SourceType: pdf50tawi.File,
				Value:      "cmd/demo-cli/demo-signature-1280x720-rectangle.png",
			},
			CompanySealImage: pdf50tawi.Image{
				SourceType: pdf50tawi.File,
				Value:      "cmd/demo-cli/demo-logo-1024x1024-square.png",
			},
			DateOfIssuance: pdf50tawi.DateOfIssuance{
				Day:   "22",
				Month: "มกราคม",
				Year:  "2568",
			},
		},
	}
}
