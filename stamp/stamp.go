package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/anuchito/pdf50tawi/fonts"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/color"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func main() {
	inputPDF := flag.String("input", "tax50tawi.pdf", "Input PDF file")
	outputPDF := flag.String("output", "tax50tawi-stamped.pdf", "Output PDF file")
	signature := flag.String("signature", "demo-signature-1024x278.png", "Signature png image file")
	flag.Parse()

	// Install fonts from embedded filesystem
	fmt.Println("Installing fonts...")
	if err := fonts.InstallFonts(); err != nil {
		log.Fatalf("Error installing fonts: %v", err)
	}

	// Add text stamp
	if err := addTextStamp(*inputPDF, *outputPDF, *signature); err != nil {
		log.Fatalf("Error adding text stamp: %v", err)
	}

	fmt.Println("Successfully processed PDF with Thai text stamp")
}

// TextStampConfig holds configuration for a text watermark
type TextStampConfig struct {
	Text     string
	Dx       float64
	Dy       float64
	FontSize int
	FontName string
}

func alignLeft() *types.HAlignment {
	a := types.AlignLeft
	return &a
}

// applyTextWatermark applies a text watermark with the given configuration
func applyTextWatermark(pdfCtx *model.Context, config TextStampConfig) error {
	wm, err := pdfcpu.ParseTextWatermarkDetails(config.Text, "", true, 1)
	if err != nil {
		return err
	}

	font := "THSarabunNew"
	if config.FontName != "" {
		font = config.FontName
	}

	wm.FillColor = color.Black
	wm.Dy = config.Dy
	wm.Dx = config.Dx
	wm.Diagonal = 0
	wm.Rotation = 0
	wm.Scale = 1
	wm.ScaleAbs = true
	wm.FontName = font
	wm.FontSize = config.FontSize
	wm.OnTop = true
	wm.HAlign = alignLeft()
	wm.Pos = types.TopLeft

	return api.WatermarkContext(pdfCtx, nil, wm)
}

func addTextStamp(inputPDF, outputPDF, signature string) error {
	// Check if input PDF exists
	if _, err := os.Stat(inputPDF); os.IsNotExist(err) {
		return fmt.Errorf("input PDF not found: %s", inputPDF)
	}

	pdfCtx, err := api.ReadContextFile(inputPDF)

	if err != nil {
		return err
	}

	wm1, err := pdfcpu.ParseTextWatermarkDetails("abcdef", "", true, 1)
	if err != nil {
		return err
	}

	wm1.FillColor = color.Black
	wm1.Dy = -390
	wm1.Dx = 375
	wm1.Diagonal = 0
	wm1.Rotation = 0
	wm1.Scale = 1
	wm1.ScaleAbs = true
	wm1.FontName = "THSarabunNew"
	wm1.FontSize = 14
	wm1.OnTop = true
	wm1.HAlign = alignLeft()
	wm1.Pos = types.TopLeft

	api.WatermarkContext(pdfCtx, nil, wm1)
	wm11, err := pdfcpu.ParseTextWatermarkDetails("abcdefABC", "", true, 1)
	if err != nil {
		return err
	}

	wm11.FillColor = color.Black
	wm11.Dy = -400
	wm11.Dx = 375
	wm11.Diagonal = 0
	wm11.Rotation = 0
	wm11.Scale = 1
	wm11.ScaleAbs = true
	wm11.FontName = "THSarabunNew"
	wm11.FontSize = 14
	wm11.OnTop = true
	wm11.HAlign = alignLeft()
	wm11.Pos = types.TopLeft

	api.WatermarkContext(pdfCtx, nil, wm11)

	wm2, err := pdfcpu.ParseTextWatermarkDetails("ABCDEF", "", true, 1)
	if err != nil {
		return err
	}

	wm2.FillColor = color.Black
	wm2.Dy = -430
	wm2.Dx = -190
	wm2.Diagonal = 0
	wm2.Rotation = 0
	wm2.Scale = 1
	wm2.ScaleAbs = true
	wm2.FontName = "THSarabunNew"
	wm2.FontSize = 14
	wm2.OnTop = true
	wm2.HAlign = alignLeft()
	wm2.Pos = types.TopRight

	api.WatermarkContext(pdfCtx, nil, wm2)

	wm22, err := pdfcpu.ParseTextWatermarkDetails("ABCDEFabc", "", true, 1)
	if err != nil {
		return err
	}

	wm22.FillColor = color.Black
	wm22.Dy = -440
	wm22.Dx = -190
	wm22.Diagonal = 0
	wm22.Rotation = 0
	wm22.Scale = 1
	wm22.ScaleAbs = true
	wm22.FontName = "THSarabunNew"
	wm22.FontSize = 14
	wm22.OnTop = true
	wm22.HAlign = alignLeft()
	wm22.Pos = types.TopRight

	api.WatermarkContext(pdfCtx, nil, wm22)

	// Define text stamps configuration with demo data - adjusted for Form 50 ทวิ layout
	textStamps := []TextStampConfig{
		// Document Details (top right)
		{Text: "001", Dx: 525, Dy: -48, FontSize: 14},      // bookNumber
		{Text: "2024-001", Dx: 525, Dy: -62, FontSize: 14}, // documentNumber

		// Payer Information (ผู้จ่ายเงิน)
		{Text: "1 2 3 4 5 6 7 8 9 0 1 2 3", Dx: 378, Dy: -81, FontSize: 16},                           // payer taxId (13 digits with spaces)
		{Text: "1234567890", Dx: 350, Dy: -95, FontSize: 14},                                          // payer taxId10Digit
		{Text: "บริษัท ตัวอย่าง จำกัด", Dx: 60, Dy: -97, FontSize: 14},                                // payer name
		{Text: "123 ถนนสุขุมวิท แขวงคลองตัน เขตวัฒนา กรุงเทพฯ 10110", Dx: 60, Dy: -115, FontSize: 12}, // payer address

		// Payee Information (ผู้ถูกหักภาษี ณ ที่จ่าย)
		{Text: "9 8 7 6 5 4 3 2 1 0 9 8 7", Dx: 378, Dy: -160, FontSize: 16}, // payee taxId (13 digits with spaces)
		{Text: "9876543210", Dx: 350, Dy: -175, FontSize: 14},                // payee taxId10Digit
		{Text: "นางสาวสมหญิง นามสกุลยาวไหม", Dx: 60, Dy: -180, FontSize: 14}, // payee name

		// Tax Filing Reference (ลำดับที่)
		{Text: "001", Dx: 100, Dy: -210, FontSize: 14}, // sequenceNumber

		{Text: "m✔x✓r✓/", Dx: 280, Dy: -240, FontSize: 16, FontName: "Times-Roman"},

		// Income Details - Row 1 (เงินเดือน ค่าจาง)
		{Text: "01/01/2024", Dx: 350, Dy: -260, FontSize: 14}, // datePaid
		{Text: "500,000.00", Dx: 450, Dy: -260, FontSize: 14}, // amountPaid
		{Text: "25,000.00", Dx: 520, Dy: -260, FontSize: 14},  // taxWithheld

		// Income Details - Row 2 (ค่าธรรมเนียม ค่านายหน้า)
		{Text: "01/06/2024", Dx: 350, Dy: -285, FontSize: 14}, // datePaid
		{Text: "50,000.00", Dx: 450, Dy: -285, FontSize: 14},  // amountPaid
		{Text: "1,500.00", Dx: 520, Dy: -285, FontSize: 14},   // taxWithheld

		// Income Details - Row 3 (ค่าแห่งลิขสิทธิ์)
		{Text: "15/08/2024", Dx: 350, Dy: -310, FontSize: 14}, // datePaid
		{Text: "30,000.00", Dx: 450, Dy: -310, FontSize: 14},  // amountPaid
		{Text: "1,500.00", Dx: 520, Dy: -310, FontSize: 14},   // taxWithheld

		// Income Details - Row 4 (ดอกเบี้ย เงินปันผล)
		{Text: "01/12/2024", Dx: 350, Dy: -335, FontSize: 14}, // datePaid
		{Text: "25,000.00", Dx: 450, Dy: -335, FontSize: 14},  // amountPaid
		{Text: "750.00", Dx: 520, Dy: -335, FontSize: 14},     // taxWithheld

		// Income Details - Row 5 (อื่นๆ)
		{Text: "31/12/2024", Dx: 350, Dy: -360, FontSize: 14}, // datePaid
		{Text: "100,000.00", Dx: 450, Dy: -360, FontSize: 14}, // amountPaid
		{Text: "5,000.00", Dx: 520, Dy: -360, FontSize: 14},   // taxWithheld

		// Totals (รวม)
		{Text: "705,000.00", Dx: 450, Dy: -400, FontSize: 14}, // totalAmountPaid
		{Text: "32,750.00", Dx: 520, Dy: -400, FontSize: 14},  // totalTaxWithheld
		{Text: "สามหมื่นสองพันเจ็ดร้อยห้าสิบบาทถ้วน", Dx: 60, Dy: -420, FontSize: 14}, // totalTaxWithheldInWords

		// Other Payments (เงินที่จ่ายเข้ากองทุน)
		{Text: "12,000.00", Dx: 180, Dy: -460, FontSize: 14}, // governmentPensionFund
		{Text: "8,000.00", Dx: 320, Dy: -460, FontSize: 14},  // socialSecurityFund
		{Text: "15,000.00", Dx: 460, Dy: -460, FontSize: 14}, // providentFund

		// Certification (ลงชื่อ ผู้จ่ายเงิน และวันที่)
		{Text: "01/01/2025", Dx: 400, Dy: -550, FontSize: 14}, // dateOfIssuance
	}

	// Apply all text stamps
	for _, stamp := range textStamps {
		if err := applyTextWatermark(pdfCtx, stamp); err != nil {
			return err
		}
	}

	// Define image stamps configuration
	// "pos:bl, off:386 72, scale:0.2 abs, op:1, rot:0"

	// readfile from signature to base64
	file, err := os.Open(signature)
	if err != nil {
		return err
	}
	defer file.Close()

	// then wrap to reader
	// reader := base64.NewDecoder(base64.StdEncoding, b64)

	// then parse to watermark ImageWatermarkForReader
	// wm, err := api.ImageWatermarkForReader(file, "", true, 1)
	wm, err := pdfcpu.ParseImageWatermarkDetails(signature, "", true, 1)
	if err != nil {
		return err
	}

	wm.Dy = 82
	wm.Dx = 360
	wm.Scale = 0.08
	wm.ScaleAbs = true
	wm.Opacity = 1.0
	wm.Diagonal = 0
	wm.Rotation = 0
	wm.HAlign = alignLeft()
	wm.OnTop = false
	wm.Pos = types.BottomLeft

	if err := api.WatermarkContext(pdfCtx, nil, wm); err != nil {
		return err
	}

	return api.WriteContextFile(pdfCtx, outputPDF)
}

// TODO List:
// [ ] define key value name for each
// [ ] stamp image with
