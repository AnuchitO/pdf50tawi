package pdf50tawi

import (
	"io"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/color"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/validate"
)

// TextStmap holds configuration for a text watermark
type TextStmap struct {
	Text     string
	Dx       float64
	Dy       float64
	FontSize int
	FontName string
	Position types.Anchor
}

// applyTextWatermark applies a text watermark with the given configuration
func applyTextWatermark(pdfCtx *model.Context, config TextStmap) error {
	wm, err := TextWatermark(config)
	if err != nil {
		return err
	}
	return api.WatermarkContext(pdfCtx, nil, wm)
}

func TextWatermark(config TextStmap) (*model.Watermark, error) {
	wm, err := pdfcpu.ParseTextWatermarkDetails(config.Text, "", true, 1)
	if err != nil {
		return nil, err
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

	return wm, nil
}

type ImageStamp struct {
	Reader  io.Reader
	Pos     types.Anchor
	Dx      float64
	Dy      float64
	Scale   float64
	Opacity float64
	OnTop   bool
}

func applyImageWatermark(pdfCtx *model.Context, config ImageStamp) error {
	wm, err := ImageWatermark(config)
	if err != nil {
		return err
	}

	return api.WatermarkContext(pdfCtx, nil, wm)
}

// Stamp Image take reader
func ImageWatermark(config ImageStamp) (*model.Watermark, error) {
	wm, err := api.ImageWatermarkForReader(config.Reader, "", true, false, types.POINTS)
	if err != nil {
		return nil, err
	}
	wm.Dy = config.Dy
	wm.Dx = config.Dx
	wm.Scale = config.Scale
	wm.ScaleAbs = true
	wm.Opacity = config.Opacity
	wm.Diagonal = 0
	wm.Rotation = 0
	wm.OnTop = config.OnTop
	wm.Pos = config.Pos
	return wm, nil
}

// positionTaxID13Digits creates individual text stamps for each digit of a tax ID
func positionTaxID13Digits(taxID string, dy float64, fontSize int) []TextStmap {
	digits := strings.ReplaceAll(taxID, " ", "")

	// X positions for 13-digit tax ID (with spacing to align position on each box form)
	xPositions := []float64{378, 396, 408, 420, 432, 450, 463, 474, 486, 498, 517, 529, 548}

	return position(digits, fontSize, dy, xPositions)
}

// positionTaxID10Digits creates individual text stamps for each digit of a tax ID
func positionTaxID10Digits(taxID string, dy float64, fontSize int) []TextStmap {
	digits := strings.ReplaceAll(taxID, " ", "")

	// X positions for 10-digit tax ID (with spacing to align position on each box form)
	xPositions := []float64{422, 440, 452, 464, 476, 494, 506, 518, 530, 548}

	return position(digits, fontSize, dy, xPositions)
}

func position(digits string, fontSize int, dy float64, xPositions []float64) []TextStmap {
	var stamps []TextStmap
	for i, digit := range digits {
		if i < len(xPositions) {
			stamps = append(stamps, TextStmap{
				Text:     string(digit),
				Dx:       xPositions[i],
				Dy:       dy,
				FontSize: fontSize,
				Position: types.TopLeft,
			})
		}
	}
	return stamps
}

func tick(pnd bool) string {
	if pnd {
		return string(rune(52)) // "✔" ZapfDingbats is one of the 14 standard "Core Fonts" defined in the original PDF specification.
	}
	return " "
}

func checkmarkStamp(isSet bool, dx float64, dy float64) TextStmap {
	return TextStmap{
		Text:     tick(isSet),
		Dx:       dx,
		Dy:       dy,
		FontSize: 10,
		Position: types.TopLeft,
		FontName: "ZapfDingbats",
	}
}

// convert data from TaxInfo to TextStampConfig
func TextStampsFromTaxInfo(tax TaxInfo) []TextStmap {

	// Payer Information (ผู้จ่ายเงิน)
	payer := []TextStmap{
		{Text: tax.Payer.Name, Dx: 58, Dy: -98, FontSize: 14, Position: types.TopLeft},
		{Text: tax.Payer.Address, Dx: 62, Dy: -124, FontSize: 12, Position: types.TopLeft},
	}
	payer = append(payer, positionTaxID13Digits(tax.Payer.TaxID, -81, 16)...)
	payer = append(payer, positionTaxID10Digits(tax.Payer.TaxID10Digit, -98, 16)...)

	// Payee Information (ผู้ถูกหักภาษี ณ ที่จ่าย)
	payee := []TextStmap{
		{Text: tax.Payee.Name, Dx: 58, Dy: -170, FontSize: 14, Position: types.TopLeft},
		{Text: tax.Payee.Address, Dx: 62, Dy: -199, FontSize: 12, Position: types.TopLeft},
	}
	payee = append(payee, positionTaxID13Digits(tax.Payee.TaxID, -150, 16)...)
	payee = append(payee, positionTaxID10Digits(tax.Payee.TaxID10Digit, -169, 16)...)
	payee = append(payee, []TextStmap{
		// Tax Filing Reference (ลำดับที่)
		{Text: tax.Payee.SequenceNumber, Dx: -190, Dy: -225, FontSize: 14, Position: types.TopCenter},

		checkmarkStamp(tax.Payee.Pnd_1a, 211.5, -230),
		checkmarkStamp(tax.Payee.Pnd_1aSpecial, 289, -230),
		checkmarkStamp(tax.Payee.Pnd_2, 397, -230),
		checkmarkStamp(tax.Payee.Pnd_2a, 211.5, -248),
		checkmarkStamp(tax.Payee.Pnd_3, 474, -230),
		checkmarkStamp(tax.Payee.Pnd_3a, 289, -248),
		checkmarkStamp(tax.Payee.Pnd_53, 397, -248),
	}...)

	// Define text stamps configuration with demo data - adjusted for Form 50 ทวิ layout
	textStamps := []TextStmap{
		// Document Details (top right)
		{Text: tax.DocumentDetails.BookNumber, Dx: 525, Dy: -48, FontSize: 14, Position: types.TopLeft},
		{Text: tax.DocumentDetails.DocumentNumber, Dx: 525, Dy: -62, FontSize: 14, Position: types.TopLeft},

		// Position: Bottom Right
		// Income Details - Row 1 (เงินเดือน ค่าจาง)
		{Text: tax.Income40_1.DatePaid, Dx: 69, Dy: 530, FontSize: 14, Position: types.BottomCenter},
		{Text: tax.Income40_1.AmountPaid, Dx: -109.5, Dy: 530, FontSize: 14, Position: types.BottomRight},
		{Text: tax.Income40_1.TaxWithheld, Dx: -38, Dy: 530, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 2 (ค่าธรรมเนียม ค่านายหน้า)
		{Text: tax.Income40_2.DatePaid, Dx: 69, Dy: 516, FontSize: 14, Position: types.BottomCenter},
		{Text: tax.Income40_2.AmountPaid, Dx: -109.5, Dy: 516, FontSize: 14, Position: types.BottomRight},
		{Text: tax.Income40_2.TaxWithheld, Dx: -38, Dy: 516, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 3 (ค่าแห่งลิขสิทธิ์)
		{Text: tax.Income40_3.DatePaid, Dx: 69, Dy: 502, FontSize: 14, Position: types.BottomCenter},
		{Text: tax.Income40_3.AmountPaid, Dx: -109.5, Dy: 502, FontSize: 14, Position: types.BottomRight},
		{Text: tax.Income40_3.TaxWithheld, Dx: -38, Dy: 502, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 4 (ดอกเบี้ย เงินปันผล) 40 (4) (ก)
		{Text: tax.Income40_4A.DatePaid, Dx: 69, Dy: 488, FontSize: 14, Position: types.BottomCenter},
		{Text: tax.Income40_4A.AmountPaid, Dx: -109.5, Dy: 488, FontSize: 14, Position: types.BottomRight},
		{Text: tax.Income40_4A.TaxWithheld, Dx: -38, Dy: 488, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 4 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (1) (1.1)
		{Text: tax.Income40_4B_1_1.DatePaid, Dx: 69, Dy: 429, FontSize: 14, Position: types.BottomCenter},
		{Text: tax.Income40_4B_1_1.AmountPaid, Dx: -109.5, Dy: 429, FontSize: 14, Position: types.BottomRight},
		{Text: tax.Income40_4B_1_1.TaxWithheld, Dx: -38, Dy: 429, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 5 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (1) (1.2)
		{Text: tax.Income40_4B_1_2.DatePaid, Dx: 69, Dy: 414, FontSize: 14, Position: types.BottomCenter},
		{Text: tax.Income40_4B_1_2.AmountPaid, Dx: -109.5, Dy: 414, FontSize: 14, Position: types.BottomRight},
		{Text: tax.Income40_4B_1_2.TaxWithheld, Dx: -38, Dy: 414, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 5 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (1) (1.3)
		{Text: tax.Income40_4B_1_3.DatePaid, Dx: 69, Dy: 400, FontSize: 14, Position: types.BottomCenter},
		{Text: tax.Income40_4B_1_3.AmountPaid, Dx: -109.5, Dy: 400, FontSize: 14, Position: types.BottomRight},
		{Text: tax.Income40_4B_1_3.TaxWithheld, Dx: -38, Dy: 400, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 6 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (1) (1.3)
		{Text: tax.Income40_4B_1_4_Rate, Dx: -116, Dy: 384, FontSize: 14, Position: types.BottomCenter},
		{Text: tax.Income40_4B_1_4.DatePaid, Dx: 69, Dy: 386, FontSize: 14, Position: types.BottomCenter},
		{Text: tax.Income40_4B_1_4.AmountPaid, Dx: -109.5, Dy: 386, FontSize: 14, Position: types.BottomRight},
		{Text: tax.Income40_4B_1_4.TaxWithheld, Dx: -38, Dy: 386, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 6 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (2) (2.1)
		{Text: tax.Income40_4B_2_1.DatePaid, Dx: 69, Dy: 356, FontSize: 14, Position: types.BottomCenter},
		{Text: tax.Income40_4B_2_1.AmountPaid, Dx: -109.5, Dy: 356, FontSize: 14, Position: types.BottomRight},
		{Text: tax.Income40_4B_2_1.TaxWithheld, Dx: -38, Dy: 356, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 6 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (2) (2.2)
		{Text: tax.Income40_4B_2_2.DatePaid, Dx: 69, Dy: 327, FontSize: 14, Position: types.BottomCenter},
		{Text: tax.Income40_4B_2_2.AmountPaid, Dx: -109.5, Dy: 327, FontSize: 14, Position: types.BottomRight},
		{Text: tax.Income40_4B_2_2.TaxWithheld, Dx: -38, Dy: 327, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 6 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (2) (2.3)
		{Text: tax.Income40_4B_2_3.DatePaid, Dx: 69, Dy: 298, FontSize: 14, Position: types.BottomCenter},
		{Text: tax.Income40_4B_2_3.AmountPaid, Dx: -109.5, Dy: 298, FontSize: 14, Position: types.BottomRight},
		{Text: tax.Income40_4B_2_3.TaxWithheld, Dx: -38, Dy: 298, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 6 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (2) (2.4)
		{Text: tax.Income40_4B_2_4.DatePaid, Dx: 69, Dy: 282, FontSize: 14, Position: types.BottomCenter},
		{Text: tax.Income40_4B_2_4.AmountPaid, Dx: -109.5, Dy: 282, FontSize: 14, Position: types.BottomRight},
		{Text: tax.Income40_4B_2_4.TaxWithheld, Dx: -38, Dy: 282, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 6 (ดอกเบี้ย เงินปันผล) 40 (4) (ข) (2) (2.5)
		{Text: tax.Income40_4B_2_5_Note, Dx: 150, Dy: 269, FontSize: 12, Position: types.BottomLeft},
		{Text: tax.Income40_4B_2_5.DatePaid, Dx: 69, Dy: 268, FontSize: 14, Position: types.BottomCenter},
		{Text: tax.Income40_4B_2_5.AmountPaid, Dx: -109.5, Dy: 268, FontSize: 14, Position: types.BottomRight},
		{Text: tax.Income40_4B_2_5.TaxWithheld, Dx: -38, Dy: 268, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 7 5. การจ่ายเงินได้ที่ต้องหักภาษี ณ ที่จ่าย
		{Text: tax.Income5.DatePaid, Dx: 69, Dy: 212, FontSize: 14, Position: types.BottomCenter},
		{Text: tax.Income5.AmountPaid, Dx: -109.5, Dy: 212, FontSize: 14, Position: types.BottomRight},
		{Text: tax.Income5.TaxWithheld, Dx: -38, Dy: 212, FontSize: 14, Position: types.BottomRight},

		// Income Details - Row 8 6. อื่น ๆ (ระบุ)
		{Text: tax.Income6_Note, Dx: 102, Dy: 196, FontSize: 12, Position: types.BottomLeft},
		{Text: tax.Income6.DatePaid, Dx: 69, Dy: 198, FontSize: 14, Position: types.BottomCenter},
		{Text: tax.Income6.AmountPaid, Dx: -109.5, Dy: 198, FontSize: 14, Position: types.BottomRight},
		{Text: tax.Income6.TaxWithheld, Dx: -38, Dy: 198, FontSize: 14, Position: types.BottomRight},
		//
		// Totals (รวม)
		{Text: tax.Totals.TotalAmountPaid, Dx: -109.5, Dy: 176, FontSize: 14, Position: types.BottomRight},
		{Text: tax.Totals.TotalTaxWithheld, Dx: -38, Dy: 176, FontSize: 14, Position: types.BottomRight},
		{Text: tax.Totals.TotalTaxWithheldInWords, Dx: 190, Dy: 156, FontSize: 14, Position: types.BottomLeft},

		// Other Payments (เงินที่จ่ายเข้ากองทุน)
		{Text: tax.OtherPayments.GovernmentPensionFund, Dx: -318, Dy: 139, FontSize: 12, Position: types.BottomRight},
		{Text: tax.OtherPayments.SocialSecurityFund, Dx: -190, Dy: 139, FontSize: 12, Position: types.BottomRight},
		{Text: tax.OtherPayments.ProvidentFund, Dx: -54, Dy: 139, FontSize: 12, Position: types.BottomRight},

		// Withholding Type
		checkmarkStamp(tax.WithholdingType.WithholdingTax, 85, -712),
		checkmarkStamp(tax.WithholdingType.Forever, 178, -712),
		checkmarkStamp(tax.WithholdingType.OneTime, 285.5, -712),
		checkmarkStamp(tax.WithholdingType.Other, 396, -712),
		{Text: tax.WithholdingType.OtherDetails, Dx: 470, Dy: 117, FontSize: 12, Position: types.BottomLeft},

		// Certification (ลงชื่อ ผู้จ่ายเงิน และวันที่)
		{Text: tax.Certification.DateOfIssuance.Day, Dx: 52, Dy: 70, FontSize: 14, Position: types.BottomCenter},
		{Text: tax.Certification.DateOfIssuance.Month, Dx: 99, Dy: 70, FontSize: 14, Position: types.BottomCenter},
		{Text: tax.Certification.DateOfIssuance.Year, Dx: 152, Dy: 70, FontSize: 14, Position: types.BottomCenter},
	}

	textStamps = append(textStamps, payer...)
	textStamps = append(textStamps, payee...)

	return textStamps
}

func CertificateImageStamps(sign io.Reader, logo io.Reader) []ImageStamp {
	return []ImageStamp{
		{Reader: sign, Pos: types.BottomCenter, Dx: 105, Dy: 84, Scale: 0.08, Opacity: 1, OnTop: true},
		{Reader: logo, Pos: types.BottomCenter, Dx: 230, Dy: 64, Scale: 0.08, Opacity: 1, OnTop: true},
	}
}

func ReadContext(inFile io.ReadSeeker) (*model.Context, error) {
	ctx, err := api.ReadContext(inFile, model.NewDefaultConfiguration())
	if err != nil {
		return nil, err
	}

	if ctx.Conf.Version != model.VersionStr {
		model.CheckConfigVersion(ctx.Conf.Version)
	}

	return ctx, validate.XRefTable(ctx)
}

// BuildStampedContext // take inputPDF and  return *model.Context
func BuildStampedContext(inputPDF io.ReadSeeker, textStamps []TextStmap, imageStamps []ImageStamp) (*model.Context, error) {
	// Ensure fonts are installed before any watermarking occurs
	if err := InstallFonts(); err != nil {
		return nil, err
	}
	ctx, err := ReadContext(inputPDF)
	if err != nil {
		return nil, err
	}

	for _, stamp := range textStamps {
		if err := applyTextWatermark(ctx, stamp); err != nil {
			return nil, err
		}
	}

	for _, stamp := range imageStamps {
		if err := applyImageWatermark(ctx, stamp); err != nil {
			return nil, err
		}
	}

	return ctx, nil
}

func WriteStampedPDF(ctx *model.Context, outputPDF io.Writer) error {
	return api.WriteContext(ctx, outputPDF)
}

func IssueWHTCertificatePDF(outputPDF io.Writer, taxInfo TaxInfo, sign io.Reader, logo io.Reader) error {
	inputPDF, err := Tax50tawiPDFTemplate()
	if err != nil {
		return err
	}
	return WHTCertificatePDF(inputPDF, outputPDF, taxInfo, sign, logo)
}

// WHTCertificatePDF
func WHTCertificatePDF(inputPDF io.ReadSeeker, outputPDF io.Writer, taxInfo TaxInfo, sign io.Reader, logo io.Reader) error {
	images := CertificateImageStamps(sign, logo)
	texts := TextStampsFromTaxInfo(taxInfo)

	ctx, err := BuildStampedContext(inputPDF, texts, images)
	if err != nil {
		return err
	}

	return WriteStampedPDF(ctx, outputPDF)
}
