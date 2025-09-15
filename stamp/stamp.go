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
	logo := flag.String("logo", "demo-logo-410x361.png", "Logo png image file")
	flag.Parse()

	// Install fonts from embedded filesystem
	fmt.Println("Installing fonts...")
	if err := fonts.InstallFonts(); err != nil {
		log.Fatalf("Error installing fonts: %v", err)
	}

	// Add text stamp
	if err := addTextStamp(*inputPDF, *outputPDF, *signature, *logo); err != nil {
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
	Position types.Anchor
}

func alignLeft() *types.HAlignment {
	a := types.AlignRight
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
	// wm.HAlign = alignLeft()
	wm.Pos = config.Position

	return api.WatermarkContext(pdfCtx, nil, wm)
}

func addTextStamp(inputPDF, outputPDF, signature, logo string) error {
	// Check if input PDF exists
	if _, err := os.Stat(inputPDF); os.IsNotExist(err) {
		return fmt.Errorf("input PDF not found: %s", inputPDF)
	}

	pdfCtx, err := api.ReadContextFile(inputPDF)

	if err != nil {
		return err
	}

	// wm1, err := pdfcpu.ParseTextWatermarkDetails("abcdef", "", true, 1)
	// if err != nil {
	// 	return err
	// }

	// wm1.FillColor = color.Black
	// wm1.Dy = 0
	// wm1.Dx = 100
	// wm1.Diagonal = 0
	// wm1.Rotation = 0
	// wm1.Scale = 1
	// wm1.ScaleAbs = true
	// wm1.FontName = "THSarabunNew"
	// wm1.FontSize = 14
	// wm1.OnTop = true
	// wm1.HAlign = alignLeft()
	// wm1.Pos = types.TopLeft

	// api.WatermarkContext(pdfCtx, nil, wm1)
	// wm11, err := pdfcpu.ParseTextWatermarkDetails("abcdefABC", "", true, 1)
	// if err != nil {
	// 	return err
	// }

	// wm11.FillColor = color.Black
	// wm11.Dy = -400
	// wm11.Dx = 375
	// wm11.Diagonal = 0
	// wm11.Rotation = 0
	// wm11.Scale = 1
	// wm11.ScaleAbs = true
	// wm11.FontName = "THSarabunNew"
	// wm11.FontSize = 14
	// wm11.OnTop = true
	// wm11.HAlign = alignLeft()
	// wm11.Pos = types.TopLeft

	// api.WatermarkContext(pdfCtx, nil, wm11)

	// wm2, err := pdfcpu.ParseTextWatermarkDetails("ABCDEF", "", true, 1)
	// if err != nil {
	// 	return err
	// }

	// wm2.FillColor = color.Black
	// wm2.Dy = -420
	// wm2.Dx = -190
	// wm2.Diagonal = 0
	// wm2.Rotation = 0
	// wm2.Scale = 1
	// wm2.ScaleAbs = true
	// wm2.FontName = "THSarabunNew"
	// wm2.FontSize = 14
	// wm2.OnTop = true
	// wm2.HAlign = alignLeft()
	// wm2.Pos = types.TopRight

	// api.WatermarkContext(pdfCtx, nil, wm2)

	// wm22, err := pdfcpu.ParseTextWatermarkDetails("ABCDEFabcERF", "", true, 1)
	// if err != nil {
	// 	return err
	// }

	// wm22.FillColor = color.Black
	// wm22.Dy = 385
	// wm22.Dx = -190
	// wm22.Diagonal = 0
	// wm22.Rotation = 0
	// wm22.Scale = 1
	// wm22.ScaleAbs = true
	// wm22.FontName = "THSarabunNew"
	// wm22.FontSize = 16
	// wm22.OnTop = true
	// wm22.HAlign = alignLeft()
	// wm22.Pos = types.BottomRight

	// api.WatermarkContext(pdfCtx, nil, wm22)

	// Define text stamps configuration with demo data - adjusted for Form 50 ทวิ layout
	textStamps := []TextStampConfig{
		// Document Details (top right)
		{Text: "001", Dx: 525, Dy: -48, FontSize: 14, Position: types.TopLeft},      // bookNumber
		{Text: "2568-001", Dx: 525, Dy: -62, FontSize: 14, Position: types.TopLeft}, // documentNumber

		// Payer Information (ผู้จ่ายเงิน)
		{Text: "1 2 3 4 5 6 7 8 9 0 1 2 3", Dx: 378, Dy: -81, FontSize: 16, Position: types.TopLeft},                           // payer taxId (13 digits with spaces)
		{Text: "บริษัท ตัวอย่าง จำกัด", Dx: 58, Dy: -98, FontSize: 14, Position: types.TopLeft},                                // payer name
		{Text: "1234567890", Dx: 422, Dy: -98, FontSize: 16, Position: types.TopLeft},                                          // payer taxId10Digit
		{Text: "123 ถนนสุขุมวิท แขวงคลองตัน เขตวัฒนา กรุงเทพฯ 10110", Dx: 62, Dy: -124, FontSize: 12, Position: types.TopLeft}, // payer address

		// Payee Information (ผู้ถูกหักภาษี ณ ที่จ่าย)
		{Text: "9 8 7 6 5 4 3 2 1 0 9 8 7", Dx: 378, Dy: -150, FontSize: 16, Position: types.TopLeft},              // payee taxId (13 digits with spaces)
		{Text: "นางสาวสมหญิง นามสกุลยาวไหม", Dx: 58, Dy: -170, FontSize: 14, Position: types.TopLeft},              // payee name
		{Text: "9876543210", Dx: 422, Dy: -168, FontSize: 16, Position: types.TopLeft},                             // payee taxId10Digit
		{Text: "555 ต.ทุ่งนา  อ.ทุ่งนา  จ.ชลบุรี  12345", Dx: 62, Dy: -199, FontSize: 12, Position: types.TopLeft}, // payee address

		// Tax Filing Reference (ลำดับที่)
		{Text: "001", Dx: 100, Dy: -225, FontSize: 14, Position: types.TopLeft}, // sequenceNumber

		// TODO: fix "✔" character
		{Text: "/", Dy: -240, Dx: 398, FontSize: 22, Position: types.TopLeft, FontName: "THSarabunNew-Bold"},
		{Text: "/", Dy: -240, Dx: 291, FontSize: 22, Position: types.TopLeft, FontName: "THSarabunNew-Bold"},
		{Text: "/", Dy: -240, Dx: 213, FontSize: 22, Position: types.TopLeft, FontName: "THSarabunNew-Bold"},
		{Text: "/", Dy: -222, Dx: 476, FontSize: 22, Position: types.TopLeft, FontName: "THSarabunNew-Bold"},
		{Text: "/", Dy: -222, Dx: 398, FontSize: 22, Position: types.TopLeft, FontName: "THSarabunNew-Bold"},
		{Text: "/", Dy: -222, Dx: 291, FontSize: 22, Position: types.TopLeft, FontName: "THSarabunNew-Bold"},
		{Text: "/", Dy: -222, Dx: 213, FontSize: 22, Position: types.TopLeft, FontName: "THSarabunNew-Bold"},

		// Position: Bottom Right
		// Income Details - Row 1 (เงินเดือน ค่าจาง)
		{Text: "01/01/2568", Dx: -205, Dy: 530, FontSize: 14, Position: types.BottomRight}, // datePaid
		{Text: "500,000.00", Dx: -110, Dy: 530, FontSize: 14, Position: types.BottomRight}, // amountPaid
		{Text: "25,000.00", Dx: -40, Dy: 530, FontSize: 14, Position: types.BottomRight},   // taxWithheld

		// Income Details - Row 2 (ค่าธรรมเนียม ค่านายหน้า)
		{Text: "02/06/2568", Dx: -205, Dy: 516, FontSize: 14, Position: types.BottomRight}, // datePaid
		{Text: "50,000.00", Dx: -110, Dy: 516, FontSize: 14, Position: types.BottomRight},  // amountPaid
		{Text: "1,500.00", Dx: -40, Dy: 516, FontSize: 14, Position: types.BottomRight},    // taxWithheld

		// Income Details - Row 3 (ค่าแห่งลิขสิทธิ์)
		{Text: "15/08/2568", Dx: -205, Dy: 502, FontSize: 14, Position: types.BottomRight}, // datePaid
		{Text: "30,000.00", Dx: -110, Dy: 502, FontSize: 14, Position: types.BottomRight},  // amountPaid
		{Text: "1,500.00", Dx: -40, Dy: 502, FontSize: 14, Position: types.BottomRight},    // taxWithheld

		// Income Details - Row 4 (ดอกเบี้ย เงินปันผล) 40 (4) (ข)
		{Text: "03/12/2568", Dx: -205, Dy: 488, FontSize: 14, Position: types.BottomRight}, // datePaid
		{Text: "25,000.00", Dx: -110, Dy: 488, FontSize: 14, Position: types.BottomRight},  // amountPaid
		{Text: "750.00", Dx: -40, Dy: 488, FontSize: 14, Position: types.BottomRight},      // taxWithheld

		// Income Details - Row 4 (ดอกเบี้ย เงินปันผล) 40 (4) (ข)
		// (1) (1.1)
		{Text: "04/12/2568", Dx: -205, Dy: 414, FontSize: 14, Position: types.BottomRight}, // datePaid
		{Text: "12,000.00", Dx: -110, Dy: 414, FontSize: 14, Position: types.BottomRight},  // amountPaid
		{Text: "750.00", Dx: -40, Dy: 414, FontSize: 14, Position: types.BottomRight},      // taxWithheld

		// Income Details - Row 5 (ดอกเบี้ย เงินปันผล) 40 (4) (ข)
		// (1) (1.2)
		{Text: "05/12/2568", Dx: -205, Dy: 400, FontSize: 14, Position: types.BottomRight}, // datePaid
		{Text: "12,000.00", Dx: -110, Dy: 400, FontSize: 14, Position: types.BottomRight},  // amountPaid
		{Text: "750.00", Dx: -40, Dy: 400, FontSize: 14, Position: types.BottomRight},      // taxWithheld

		// Income Details - Row 6 (ดอกเบี้ย เงินปันผล) 40 (4) (ข)
		// (1) (1.3)
		// TODO: fix % symbol not showing
		{Text: "12", Dx: -405, Dy: 384, FontSize: 12, Position: types.BottomRight},         // otherRate
		{Text: "06/12/2568", Dx: -205, Dy: 386, FontSize: 14, Position: types.BottomRight}, // datePaid
		{Text: "12,000.00", Dx: -110, Dy: 386, FontSize: 14, Position: types.BottomRight},  // amountPaid
		{Text: "750.00", Dx: -40, Dy: 386, FontSize: 14, Position: types.BottomRight},      // taxWithheld

		// Income Details - Row 6 (ดอกเบี้ย เงินปันผล) 40 (4) (ข)
		// (2) (2.1)
		{Text: "21/12/2568", Dx: -205, Dy: 356, FontSize: 14, Position: types.BottomRight}, // datePaid
		{Text: "100,222.00", Dx: -110, Dy: 356, FontSize: 14, Position: types.BottomRight}, // amountPaid
		{Text: "5,222.00", Dx: -40, Dy: 356, FontSize: 14, Position: types.BottomRight},    // taxWithheld

		// Income Details - Row 6 (ดอกเบี้ย เงินปันผล) 40 (4) (ข)
		// (2) (2.2)
		{Text: "22/12/2568", Dx: -205, Dy: 327, FontSize: 14, Position: types.BottomRight}, // datePaid
		{Text: "100,222.00", Dx: -110, Dy: 327, FontSize: 14, Position: types.BottomRight}, // amountPaid
		{Text: "5,222.00", Dx: -40, Dy: 327, FontSize: 14, Position: types.BottomRight},    // taxWithheld

		// Income Details - Row 6 (ดอกเบี้ย เงินปันผล) 40 (4) (ข)
		// (2) (2.3)
		{Text: "23/12/2568", Dx: -205, Dy: 298, FontSize: 14, Position: types.BottomRight}, // datePaid
		{Text: "100,233.00", Dx: -110, Dy: 298, FontSize: 14, Position: types.BottomRight}, // amountPaid
		{Text: "5,233.00", Dx: -40, Dy: 298, FontSize: 14, Position: types.BottomRight},    // taxWithheld

		// Income Details - Row 6 (ดอกเบี้ย เงินปันผล) 40 (4) (ข)
		// (2) (2.4)
		{Text: "24/12/2568", Dx: -205, Dy: 282, FontSize: 14, Position: types.BottomRight}, // datePaid
		{Text: "100,244.00", Dx: -110, Dy: 282, FontSize: 14, Position: types.BottomRight}, // amountPaid
		{Text: "5,244.00", Dx: -40, Dy: 282, FontSize: 14, Position: types.BottomRight},    // taxWithheld

		// Income Details - Row 6 (ดอกเบี้ย เงินปันผล) 40 (4) (ข)
		// (2) (2.5)
		{Text: "กำไรจากที่ไหนสักที่", Dx: 150, Dy: 268, FontSize: 12, Position: types.BottomLeft}, // otherRate
		{Text: "25/12/2568", Dx: -205, Dy: 270, FontSize: 14, Position: types.BottomRight},        // datePaid
		{Text: "100,255.00", Dx: -110, Dy: 270, FontSize: 14, Position: types.BottomRight},        // amountPaid
		{Text: "5,255.00", Dx: -40, Dy: 270, FontSize: 14, Position: types.BottomRight},           // taxWithheld

		// Income Details - Row 7 5. การจ่ายเงินได้ที่ต้องหักภาษี ณ ที่จ่าย
		{Text: "05/12/2568", Dx: -205, Dy: 248, FontSize: 14, Position: types.BottomRight}, // datePaid
		{Text: "100,555.00", Dx: -110, Dy: 248, FontSize: 14, Position: types.BottomRight}, // amountPaid
		{Text: "5,555.00", Dx: -40, Dy: 248, FontSize: 14, Position: types.BottomRight},    // taxWithheld

		// Income Details - Row 8 6. อื่น ๆ (ระบุ)
		{Text: "อื่นอื่นอื่นอื่นอื่นอื่นอื่นอื่นอื่น", Dx: 102, Dy: 196, FontSize: 12, Position: types.BottomLeft}, // otherRate
		{Text: "06/12/2568", Dx: -205, Dy: 198, FontSize: 14, Position: types.BottomRight},                         // datePaid
		{Text: "100,666.00", Dx: -110, Dy: 198, FontSize: 14, Position: types.BottomRight},                         // amountPaid
		{Text: "6,666.00", Dx: -40, Dy: 198, FontSize: 14, Position: types.BottomRight},                            // taxWithheld

		// Totals (รวม)
		{Text: "705,000.00", Dx: -110, Dy: 176, FontSize: 14, Position: types.BottomRight}, // totalAmountPaid
		{Text: "32,750.00", Dx: -40, Dy: 176, FontSize: 14, Position: types.BottomRight},   // totalTaxWithheld

		{Text: "สามหมื่นสองพันเจ็ดร้อยห้าสิบบาทถ้วน", Dx: 190, Dy: 156, FontSize: 14, Position: types.BottomLeft}, // totalTaxWithheldInWords

		// Other Payments (เงินที่จ่ายเข้ากองทุน)
		{Text: "12,000.00", Dx: -318, Dy: 139, FontSize: 12, Position: types.BottomRight}, // governmentPensionFund
		{Text: "8,000.00", Dx: -190, Dy: 139, FontSize: 12, Position: types.BottomRight},  // socialSecurityFund
		{Text: "15,000.00", Dx: -54, Dy: 139, FontSize: 12, Position: types.BottomRight},  // providentFund

		{Text: "/", Dx: 86, Dy: 110, FontSize: 22, Position: types.BottomLeft, FontName: "THSarabunNew-Bold"},
		{Text: "/", Dx: 178, Dy: 110, FontSize: 22, Position: types.BottomLeft, FontName: "THSarabunNew-Bold"},
		{Text: "/", Dx: 286, Dy: 110, FontSize: 22, Position: types.BottomLeft, FontName: "THSarabunNew-Bold"},
		{Text: "/", Dx: 397, Dy: 110, FontSize: 22, Position: types.BottomLeft, FontName: "THSarabunNew-Bold"},
		{Text: "อื่นอื่นอื่นอื่นอื่นอื่นอื่นอื่น", Dx: 470, Dy: 117, FontSize: 12, Position: types.BottomLeft},

		// Certification (ลงชื่อ ผู้จ่ายเงิน และวันที่)
		{Text: "99/09/2568", Dx: 370, Dy: 70, FontSize: 14, Position: types.BottomLeft}, // dateOfIssuance
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

	// Logo
	wmLogo, err := pdfcpu.ParseImageWatermarkDetails(logo, "", true, 1)
	if err != nil {
		return err
	}

	wmLogo.Dy = 64
	wmLogo.Dx = 511
	wmLogo.Scale = 0.08
	wmLogo.ScaleAbs = true
	wmLogo.Opacity = 1.0
	wmLogo.Diagonal = 0
	wmLogo.Rotation = 0
	wmLogo.HAlign = alignLeft()
	wmLogo.OnTop = true
	wmLogo.Pos = types.BottomLeft

	if err := api.WatermarkContext(pdfCtx, nil, wmLogo); err != nil {
		return err
	}

	return api.WriteContextFile(pdfCtx, outputPDF)
}

// TODO List:
// [ ] define key value name for each
// [ ] stamp image with
