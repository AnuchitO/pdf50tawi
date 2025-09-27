package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

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

func alignCenter() *types.HAlignment {
	a := types.AlignCenter
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
	wm.Pos = config.Position

	return api.WatermarkContext(pdfCtx, nil, wm)
}

// positionTaxID13Digits creates individual text stamps for each digit of a tax ID
func positionTaxID13Digits(taxID string, startX, y float64, fontSize int) []TextStampConfig {
	var stamps []TextStampConfig
	digits := strings.ReplaceAll(taxID, " ", "") // Remove spaces

	// X positions for 13-digit tax ID (with spacing for readability)
	xPositions := []float64{378, 396, 408, 420, 432, 450, 463, 474, 486, 498, 517, 529, 548}

	for i, digit := range digits {
		if i < len(xPositions) {
			stamps = append(stamps, TextStampConfig{
				Text:     string(digit),
				Dx:       xPositions[i],
				Dy:       y,
				FontSize: fontSize,
				Position: types.TopLeft,
			})
		}
	}
	return stamps
}

// positionTaxID10Digits creates individual text stamps for each digit of a tax ID
func positionTaxID10Digits(taxID string, startX, y float64, fontSize int) []TextStampConfig {
	var stamps []TextStampConfig
	digits := strings.ReplaceAll(taxID, " ", "") // Remove spaces

	// X positions for 10-digit tax ID (with spacing for readability)
	xPositions := []float64{422, 440, 452, 464, 476, 494, 506, 518, 530, 548}

	for i, digit := range digits {
		if i < len(xPositions) {
			stamps = append(stamps, TextStampConfig{
				Text:     string(digit),
				Dx:       xPositions[i],
				Dy:       y,
				FontSize: fontSize,
				Position: types.TopLeft,
			})
		}
	}
	return stamps
}

func tick(pnd bool) string {
	if pnd {
		return "/" // TODO: fix "✔" character
	}
	return " "
}

// convert data from Payload to TextStampConfig
func convertPayloadToTextStampConfig(payload Payload) []TextStampConfig {

	payer := []TextStampConfig{
		// Payer Information (ผู้จ่ายเงิน)
		{Text: payload.Payer.Name, Dx: 58, Dy: -98, FontSize: 14, Position: types.TopLeft},
		{Text: payload.Payer.Address, Dx: 62, Dy: -124, FontSize: 12, Position: types.TopLeft},
	}
	payer = append(payer, positionTaxID13Digits(payload.Payer.TaxID, 378, -81, 16)...)
	payer = append(payer, positionTaxID10Digits(payload.Payer.TaxID10Digit, 422, -98, 16)...)

	// Payee Information (ผู้ถูกหักภาษี ณ ที่จ่าย)
	payee := append([]TextStampConfig{}, positionTaxID13Digits(payload.Payee.TaxID, 378, -150, 16)...)
	payee = append(payee, positionTaxID10Digits(payload.Payee.TaxID10Digit, 422, -169, 16)...)
	payee = append(payee, []TextStampConfig{
		{Text: payload.Payee.Name, Dx: 58, Dy: -170, FontSize: 14, Position: types.TopLeft},
		{Text: payload.Payee.Address, Dx: 62, Dy: -199, FontSize: 12, Position: types.TopLeft},
		// Tax Filing Reference (ลำดับที่)
		{Text: payload.Payee.SequenceNumber, Dx: -190, Dy: -225, FontSize: 14, Position: types.TopCenter},

		// TODO: fix "✔" character
		{Text: tick(payload.Payee.Pnd_1a), Dy: -222, Dx: 213, FontSize: 22, Position: types.TopLeft, FontName: "THSarabunNew-Bold"},
		{Text: tick(payload.Payee.Pnd_1aSpecial), Dy: -222, Dx: 291, FontSize: 22, Position: types.TopLeft, FontName: "THSarabunNew-Bold"},
		{Text: tick(payload.Payee.Pnd_2), Dy: -222, Dx: 398, FontSize: 22, Position: types.TopLeft, FontName: "THSarabunNew-Bold"},
		{Text: tick(payload.Payee.Pnd_2a), Dy: -240, Dx: 213, FontSize: 22, Position: types.TopLeft, FontName: "THSarabunNew-Bold"},
		{Text: tick(payload.Payee.Pnd_3), Dy: -222, Dx: 476, FontSize: 22, Position: types.TopLeft, FontName: "THSarabunNew-Bold"},
		{Text: tick(payload.Payee.Pnd_3a), Dy: -240, Dx: 291, FontSize: 22, Position: types.TopLeft, FontName: "THSarabunNew-Bold"},
		{Text: tick(payload.Payee.Pnd_53), Dy: -240, Dx: 398, FontSize: 22, Position: types.TopLeft, FontName: "THSarabunNew-Bold"},
	}...)

	// Define text stamps configuration with demo data - adjusted for Form 50 ทวิ layout
	textStamps := []TextStampConfig{
		// Document Details (top right)
		{Text: payload.DocumentDetails.BookNumber, Dx: 525, Dy: -48, FontSize: 14, Position: types.TopLeft},
		{Text: payload.DocumentDetails.DocumentNumber, Dx: 525, Dy: -62, FontSize: 14, Position: types.TopLeft},

		// Position: Bottom Right
		// Income Details - Row 1 (เงินเดือน ค่าจาง)
		{Text: payload.Income40_1.DatePaid, Dx: 69, Dy: 530, FontSize: 14, Position: types.BottomCenter},
		{Text: payload.Income40_1.AmountPaid, Dx: -109.5, Dy: 530, FontSize: 14, Position: types.BottomRight},
		{Text: payload.Income40_1.TaxWithheld, Dx: -38, Dy: 530, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 2 (ค่าธรรมเนียม ค่านายหน้า)
		{Text: payload.Income40_2.DatePaid, Dx: 69, Dy: 516, FontSize: 14, Position: types.BottomCenter},
		{Text: payload.Income40_2.AmountPaid, Dx: -109.5, Dy: 516, FontSize: 14, Position: types.BottomRight},
		{Text: payload.Income40_2.TaxWithheld, Dx: -38, Dy: 516, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 3 (ค่าแห่งลิขสิทธิ์)
		{Text: payload.Income40_3.DatePaid, Dx: 69, Dy: 502, FontSize: 14, Position: types.BottomCenter},
		{Text: payload.Income40_3.AmountPaid, Dx: -109.5, Dy: 502, FontSize: 14, Position: types.BottomRight},
		{Text: payload.Income40_3.TaxWithheld, Dx: -38, Dy: 502, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 4 (ดอกเบี้ย เงินปันผล) 40 (4) (ก)
		{Text: payload.Income40_4A.DatePaid, Dx: 69, Dy: 488, FontSize: 14, Position: types.BottomCenter},
		{Text: payload.Income40_4A.AmountPaid, Dx: -109.5, Dy: 488, FontSize: 14, Position: types.BottomRight},
		{Text: payload.Income40_4A.TaxWithheld, Dx: -38, Dy: 488, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 4 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (1) (1.1)
		{Text: payload.Income40_4B_1_1.DatePaid, Dx: 69, Dy: 429, FontSize: 14, Position: types.BottomCenter},
		{Text: payload.Income40_4B_1_1.AmountPaid, Dx: -109.5, Dy: 429, FontSize: 14, Position: types.BottomRight},
		{Text: payload.Income40_4B_1_1.TaxWithheld, Dx: -38, Dy: 429, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 5 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (1) (1.2)
		{Text: payload.Income40_4B_1_2.DatePaid, Dx: 69, Dy: 414, FontSize: 14, Position: types.BottomCenter},
		{Text: payload.Income40_4B_1_2.AmountPaid, Dx: -109.5, Dy: 414, FontSize: 14, Position: types.BottomRight},
		{Text: payload.Income40_4B_1_2.TaxWithheld, Dx: -38, Dy: 414, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 5 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (1) (1.3)
		{Text: payload.Income40_4B_1_3.DatePaid, Dx: 69, Dy: 400, FontSize: 14, Position: types.BottomCenter},
		{Text: payload.Income40_4B_1_3.AmountPaid, Dx: -109.5, Dy: 400, FontSize: 14, Position: types.BottomRight},
		{Text: payload.Income40_4B_1_3.TaxWithheld, Dx: -38, Dy: 400, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 6 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (1) (1.3)
		{Text: payload.Income40_4B_1_4_Rate, Dx: -116, Dy: 384, FontSize: 14, Position: types.BottomCenter},
		{Text: payload.Income40_4B_1_4.DatePaid, Dx: 69, Dy: 386, FontSize: 14, Position: types.BottomCenter},
		{Text: payload.Income40_4B_1_4.AmountPaid, Dx: -109.5, Dy: 386, FontSize: 14, Position: types.BottomRight},
		{Text: payload.Income40_4B_1_4.TaxWithheld, Dx: -38, Dy: 386, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 6 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (2) (2.1)
		{Text: payload.Income40_4B_2_1.DatePaid, Dx: 69, Dy: 356, FontSize: 14, Position: types.BottomCenter},
		{Text: payload.Income40_4B_2_1.AmountPaid, Dx: -109.5, Dy: 356, FontSize: 14, Position: types.BottomRight},
		{Text: payload.Income40_4B_2_1.TaxWithheld, Dx: -38, Dy: 356, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 6 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (2) (2.2)
		{Text: payload.Income40_4B_2_2.DatePaid, Dx: 69, Dy: 327, FontSize: 14, Position: types.BottomCenter},
		{Text: payload.Income40_4B_2_2.AmountPaid, Dx: -109.5, Dy: 327, FontSize: 14, Position: types.BottomRight},
		{Text: payload.Income40_4B_2_2.TaxWithheld, Dx: -38, Dy: 327, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 6 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (2) (2.3)
		{Text: payload.Income40_4B_2_3.DatePaid, Dx: 69, Dy: 298, FontSize: 14, Position: types.BottomCenter},
		{Text: payload.Income40_4B_2_3.AmountPaid, Dx: -109.5, Dy: 298, FontSize: 14, Position: types.BottomRight},
		{Text: payload.Income40_4B_2_3.TaxWithheld, Dx: -38, Dy: 298, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 6 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (2) (2.4)
		{Text: payload.Income40_4B_2_4.DatePaid, Dx: 69, Dy: 282, FontSize: 14, Position: types.BottomCenter},
		{Text: payload.Income40_4B_2_4.AmountPaid, Dx: -109.5, Dy: 282, FontSize: 14, Position: types.BottomRight},
		{Text: payload.Income40_4B_2_4.TaxWithheld, Dx: -38, Dy: 282, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 6 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (2) (2.5)
		{Text: payload.Income40_4B_2_5_Note, Dx: 150, Dy: 269, FontSize: 12, Position: types.BottomLeft},
		{Text: payload.Income40_4B_2_5.DatePaid, Dx: 69, Dy: 268, FontSize: 14, Position: types.BottomCenter},
		{Text: payload.Income40_4B_2_5.AmountPaid, Dx: -109.5, Dy: 268, FontSize: 14, Position: types.BottomRight},
		{Text: payload.Income40_4B_2_5.TaxWithheld, Dx: -38, Dy: 268, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 7 5. การจ่ายเงินได้ที่ต้องหักภาษี ณ ที่จ่าย
		{Text: payload.Income40_5.DatePaid, Dx: 69, Dy: 212, FontSize: 14, Position: types.BottomCenter},
		{Text: payload.Income40_5.AmountPaid, Dx: -109.5, Dy: 212, FontSize: 14, Position: types.BottomRight},
		{Text: payload.Income40_5.TaxWithheld, Dx: -38, Dy: 212, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 8 6. อื่น ๆ (ระบุ)
		{Text: payload.Income40_6_Note, Dx: 102, Dy: 196, FontSize: 12, Position: types.BottomLeft},
		{Text: payload.Income40_6.DatePaid, Dx: 69, Dy: 198, FontSize: 14, Position: types.BottomCenter},
		{Text: payload.Income40_6.AmountPaid, Dx: -109.5, Dy: 198, FontSize: 14, Position: types.BottomRight},
		{Text: payload.Income40_6.TaxWithheld, Dx: -38, Dy: 198, FontSize: 14, Position: types.BottomRight},
		//
		// Totals (รวม)
		{Text: payload.Totals.TotalAmountPaid, Dx: -109.5, Dy: 176, FontSize: 14, Position: types.BottomRight},
		{Text: payload.Totals.TotalTaxWithheld, Dx: -38, Dy: 176, FontSize: 14, Position: types.BottomRight},
		{Text: payload.Totals.TotalTaxWithheldInWords, Dx: 190, Dy: 156, FontSize: 14, Position: types.BottomLeft},

		// Other Payments (เงินที่จ่ายเข้ากองทุน)
		{Text: payload.OtherPayments.GovernmentPensionFund, Dx: -318, Dy: 139, FontSize: 12, Position: types.BottomRight},
		{Text: payload.OtherPayments.SocialSecurityFund, Dx: -190, Dy: 139, FontSize: 12, Position: types.BottomRight},
		{Text: payload.OtherPayments.ProvidentFund, Dx: -54, Dy: 139, FontSize: 12, Position: types.BottomRight},

		{Text: tick(payload.WithholdingType.WithholdingTax), Dx: 86, Dy: 110, FontSize: 22, Position: types.BottomLeft, FontName: "THSarabunNew-Bold"},
		{Text: tick(payload.WithholdingType.Forever), Dx: 178, Dy: 110, FontSize: 22, Position: types.BottomLeft, FontName: "THSarabunNew-Bold"},
		{Text: tick(payload.WithholdingType.OneTime), Dx: 286, Dy: 110, FontSize: 22, Position: types.BottomLeft, FontName: "THSarabunNew-Bold"},
		{Text: tick(payload.WithholdingType.Other), Dx: 397, Dy: 110, FontSize: 22, Position: types.BottomLeft, FontName: "THSarabunNew-Bold"},
		{Text: payload.WithholdingType.OtherDetails, Dx: 470, Dy: 117, FontSize: 12, Position: types.BottomLeft},

		// Certification (ลงชื่อ ผู้จ่ายเงิน และวันที่)
		{Text: payload.Certification.DateOfIssuance.Day, Dx: 52, Dy: 70, FontSize: 14, Position: types.BottomCenter},
		{Text: payload.Certification.DateOfIssuance.Month, Dx: 99, Dy: 70, FontSize: 14, Position: types.BottomCenter},
		{Text: payload.Certification.DateOfIssuance.Year, Dx: 152, Dy: 70, FontSize: 14, Position: types.BottomCenter},
	}

	textStamps = append(textStamps, payer...)
	textStamps = append(textStamps, payee...)

	return textStamps
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

	textStamps := convertPayloadToTextStampConfig(DemoPayload())

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
	// TODO: support
	// 1. file on disk
	// 2. url https://example.com/image.png
	// 3. base64
	// 4. bytes

	// demo base64 payload.Certificate.PayerSignatureImage
	b64 := DemoPayload().Certification.PayerSignatureImageBase64
	// Remove data URL prefix if present
	if strings.HasPrefix(b64, "data:") {
		parts := strings.SplitN(b64, ",", 2)
		if len(parts) == 2 {
			b64 = parts[1]
		}
	}
	// Create a reader from the base64 string
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64))

	// Create image watermark from reader
	wm, err := ImageWatermark(reader, types.BottomCenter, 105, 84, 0.08, 1, false)
	if err != nil {
		return err
	}

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
// [ ] stamp image with base64
// [ ] copy- original and copy

// Stamp Image take reader
func ImageWatermark(reader io.Reader, pos types.Anchor, dx, dy float64, scale float64, opacity float64, onTop bool) (*model.Watermark, error) {
	wm, err := api.ImageWatermarkForReader(reader, "", true, false, types.POINTS)
	if err != nil {
		return nil, err
	}
	wm.Dy = dy
	wm.Dx = dx
	wm.Scale = scale
	wm.ScaleAbs = true
	wm.Opacity = opacity
	wm.Diagonal = 0
	wm.Rotation = 0
	wm.OnTop = onTop
	wm.Pos = pos
	return wm, nil
}
