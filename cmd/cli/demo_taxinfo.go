package main

import "github.com/AnuchitO/pdf50tawi"

func demoTaxInfo() pdf50tawi.TaxInfo {
	return pdf50tawi.TaxInfo{
		DocumentDetails: pdf50tawi.DocumentDetails{
			BookNumber:     "001",
			DocumentNumber: "001",
		},
		Payer: pdf50tawi.Payer{
			TaxID:        "1234567890123",
			TaxID10Digit: "1234567890",
			Name:         "บริษัท ตัวอย่าง จำกัด",
			Address:      "123 ถนนสุขุมวิท แขวงคลองตัน เขตวัฒนา กรุงเทพฯ 10110",
		},
		Payee: pdf50tawi.Payee{
			TaxID:          "3210987654321",
			TaxID10Digit:   "1234567890",
			Name:           "นางสาวสมชาย นามสกุลยาวมากไหมนะก็ไม่รู้เหมือนกัน",
			Address:        "555 ต.ทุ่งนา  อ.ทุ่งนา  จ.ชลบุรี  12345",
			SequenceNumber: "321",
			Pnd_1aSpecial:  true,
			Pnd_1a:         true,
			Pnd_2:          true,
			Pnd_3:          true,
			Pnd_2a:         true,
			Pnd_3a:         true,
			Pnd_53:         true,
		},
		Income40_1:           pdf50tawi.IncomeDetail{DatePaid: "01 มกราคม 2568", AmountPaid: "401,010.01", TaxWithheld: "12,030.30"},
		Income40_2:           pdf50tawi.IncomeDetail{DatePaid: "02 ก.พ. 2568", AmountPaid: "402,020.02", TaxWithheld: "12,060.60"},
		Income40_3:           pdf50tawi.IncomeDetail{DatePaid: "03 มี.ค. 2568", AmountPaid: "403,030.03", TaxWithheld: "12,090.90"},
		Income40_4A:          pdf50tawi.IncomeDetail{DatePaid: "04 เม.ย. 2568", AmountPaid: "404,040.04", TaxWithheld: "12,121.20"},
		Income40_4B_1_1:      pdf50tawi.IncomeDetail{DatePaid: "05 พ.ค. 2568", AmountPaid: "411,010.01", TaxWithheld: "12,330.30"},
		Income40_4B_1_2:      pdf50tawi.IncomeDetail{DatePaid: "06 มิ.ย. 2568", AmountPaid: "412,020.02", TaxWithheld: "12,360.60"},
		Income40_4B_1_3:      pdf50tawi.IncomeDetail{DatePaid: "07 ก.ค. 2568", AmountPaid: "413,030.03", TaxWithheld: "12,390.90"},
		Income40_4B_1_4_Rate: "ร้อยละ 7",
		Income40_4B_1_4:      pdf50tawi.IncomeDetail{DatePaid: "08 ส.ค. 2568", AmountPaid: "414,040.04", TaxWithheld: "12,421.20"},
		Income40_4B_2_1:      pdf50tawi.IncomeDetail{DatePaid: "09 ก.ย. 2568", AmountPaid: "421,010.01", TaxWithheld: "12,630.30"},
		Income40_4B_2_2:      pdf50tawi.IncomeDetail{DatePaid: "10 ต.ค. 2568", AmountPaid: "422,020.02", TaxWithheld: "12,660.60"},
		Income40_4B_2_3:      pdf50tawi.IncomeDetail{DatePaid: "11 พ.ย. 2568", AmountPaid: "423,030.03", TaxWithheld: "12,690.90"},
		Income40_4B_2_4:      pdf50tawi.IncomeDetail{DatePaid: "12 ธ.ค. 2568", AmountPaid: "424,040.04", TaxWithheld: "12,721.20"},
		Income40_4B_2_5_Note: "กำไรอื่นๆ",
		Income40_4B_2_5:      pdf50tawi.IncomeDetail{DatePaid: "13 ม.ค. 2568", AmountPaid: "425,050.05", TaxWithheld: "12,751.50"},
		Income5:              pdf50tawi.IncomeDetail{DatePaid: "14 ก.พ. 2568", AmountPaid: "500,010.01", TaxWithheld: "15,000.30"},
		Income6_Note:         "รายได้อื่นๆ",
		Income6:              pdf50tawi.IncomeDetail{DatePaid: "15 มี.ค. 2568", AmountPaid: "600,060.06", TaxWithheld: "18,001.80"},
		Totals: pdf50tawi.Totals{
			TotalAmountPaid:         "5,741,320.36",
			TotalTaxWithheld:        "172,239.60",
			TotalTaxWithheldInWords: "หนึ่งแสนเจ็ดหมื่นสองพันสองร้อยสามสิบเก้าบาทหกสิบสตางค์",
		},
		OtherPayments: pdf50tawi.OtherPayments{
			GovernmentPensionFund: "5,000.00",
			SocialSecurityFund:    "750.00",
			ProvidentFund:         "3,000.00",
		},
		WithholdingType: pdf50tawi.WithholdingType{
			WithholdingTax: true,
			Forever:        true,
			OneTime:        true,
			Other:          true,
			OtherDetails:   "อื่นๆ อื่นๆ อื่นๆ อื่นๆ",
		},
		Certification: pdf50tawi.Certification{
			DateOfIssuance: pdf50tawi.DateOfIssuance{
				Day:   "22",
				Month: "ธันวาคม",
				Year:  "2568",
			},
		},
	}
}
