package pdf50tawi

type TaxInfo struct {
	DocumentDetails DocumentDetails `json:"documentDetails"`
	Payer           Payer           `json:"payer"`
	Payee           Payee           `json:"payee"`

	Income40_1  IncomeDetail `json:"income40_1"`  // 1. เงินเดือน ค่าจาง เบี้ยเลี้ยง โบนัส ฯลฯ ตามมาตรา 40 (1)
	Income40_2  IncomeDetail `json:"income40_2"`  // 2. ค่าธรรมเนียม ค่านายหน้า ฯลฯ ตามมาตรา 40 (2)
	Income40_3  IncomeDetail `json:"income40_3"`  // 3. ค่าแห่งลิขสิทธิ์ ฯลฯ ตามมาตรา 40 (3)
	Income40_4A IncomeDetail `json:"income40_4A"` // 4. (ก) ดอกเบี้ย ฯลฯ ตามมาตรา 40 (4) (ก)

	// 4. (ข) เงินปันผล เงินส่วนแบ่งกำไร ฯลฯ ตามมาตรา 40 (4) (ข)
	// 4. (ข) (1) (1) กรณีผู้ได้รับเงินปันผลได้รับเครดิตภาษี โดยจ่ายจาก
	// กำไรสุทธิของกิจการที่ต้องเสียภาษีเงินได้นิติบุคคลในอัตราดังนี้
	Income40_4B_1_1      IncomeDetail `json:"income40_4B_1_1"`      // 4. (ข) (1) (1.1) อัตราร้อยละ 30 ของกำไรสุทธิ
	Income40_4B_1_2      IncomeDetail `json:"income40_4B_1_2"`      // 4. (ข) (1) (1.2) อัตราร้อยละ 25 ของกำไรสุทธิ
	Income40_4B_1_3      IncomeDetail `json:"income40_4B_1_3"`      // 4. (ข) (1) (1.3) อัตราร้อยละ 20 ของกำไรสุทธิ
	Income40_4B_1_4_Rate string       `json:"income40_4B_1_4_rate"` // 4. (ข) (1) (1.4) อัตราอื่น ๆ (ระบุ)... ของกำไรสุทธิ
	Income40_4B_1_4      IncomeDetail `json:"income40_4B_1_4"`      // 4. (ข) (1) (1.4)
	Income40_4B_2_1      IncomeDetail `json:"income40_4B_2_1"`      // 4. (ข) (2) (2.1) กำไรสุทธิของกิจการที่ได้รับยกเว้นภาษีเงินได้นิติบุคคล
	Income40_4B_2_2      IncomeDetail `json:"income40_4B_2_2"`      // 4. (ข) (2) (2.2) เงินปันผลหรือเงินส่วนแบ่งของกำไรที่ได้รับยกเว้นไม่ต้องนำมารวม คำนวณเป็นรายได้เพื่อเสียภาษีเงินได้นิติบุคคล
	Income40_4B_2_3      IncomeDetail `json:"income40_4B_2_3"`      // 4. (ข) (2) (2.3) กำไรสุทธิส่วนที่ได้หักผลขาดทุนสุทธิยกมาไม่เกิน 5 ปี ก่อนรอบระยะเวลาบัญชีปีปัจจุบัน
	Income40_4B_2_4      IncomeDetail `json:"income40_4B_2_4"`      // 4. (ข) (2) (2.4)  กำไรที่รับรู้ทางบัญชีโดยวิธีส่วนได้เสีย (equity method)
	Income40_4B_2_5_Note string       `json:"income40_4B_2_5_note"` // 4. (ข) (2) (2.5) อื่น ๆ (ระบุ)... ของกำไรสุทธิ
	Income40_4B_2_5      IncomeDetail `json:"income40_4B_2_5"`      // 4. (ข) (2) (2.5)

	Income5      IncomeDetail `json:"income5"`      // 5. การจ่ายเงินได้ที่ต้องหักภาษี ณ ที่จ่าย
	Income6      IncomeDetail `json:"income6"`      // 6. อื่น ๆ (ระบุ)
	Income6_Note string       `json:"income6_note"` // 6. อื่น ๆ (ระบุ)

	Totals          Totals          `json:"totals"`          // รวมเงิน
	TotalsInWords   string          `json:"totalsInWords"`   // รวมเงิน (ตัวอักษร)
	OtherPayments   OtherPayments   `json:"otherPayments"`   // จ่ายภาษี
	WithholdingType WithholdingType `json:"withholdingType"` // ประเภทการหักภาษี
	Certification   Certification   `json:"certification"`   // การยืนยัน
}

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
	Pnd_1a         bool   `json:"pnd_1a"`        // ภ.ง.ด. 1ก
	Pnd_1aSpecial  bool   `json:"pnd_1aSpecial"` // ภ.ง.ด. 1ก พิเศษ
	Pnd_2          bool   `json:"pnd_2"`         // ภ.ง.ด. 2
	Pnd_3          bool   `json:"pnd_3"`         // ภ.ง.ด. 3
	Pnd_2a         bool   `json:"pnd_2a"`        // ภ.ง.ด. 2ก
	Pnd_3a         bool   `json:"pnd_3a"`        // ภ.ง.ด. 3ก
	Pnd_53         bool   `json:"pnd_53"`        // ภ.ง.ด. 53
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

type WithholdingType struct {
	WithholdingTax bool   `json:"withholdingTax"` // หัก ณ ที่จ่าย
	Forever        bool   `json:"forever"`        // ออกให้ตลอดไป
	OneTime        bool   `json:"oneTime"`        // ออกให้ครั้งเดียว
	Other          bool   `json:"other"`          // อื่น ๆ
	OtherDetails   string `json:"otherDetails"`   // อื่น ๆ (ระบุ)
}

type DateOfIssuance struct {
	Day   string `json:"day"`
	Month string `json:"month"`
	Year  string `json:"year"`
}

type SourceType string

const (
	Upload SourceType = "upload"
	URL    SourceType = "url"
	File   SourceType = "file"
)

func (s SourceType) String() string {
	return string(s)
}

type Image struct {
	SourceType SourceType `json:"sourceType"` // "upload" or "url" or "file"
	Value      string     `json:"value"`
}

type Certification struct {
	PayerSignatureImage Image          `json:"payerSignatureImage"`
	CompanySealImage    Image          `json:"companySealImage"`
	DateOfIssuance      DateOfIssuance `json:"dateOfIssuance"`
}
