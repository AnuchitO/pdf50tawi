package pdf50tawi

import (
	"io"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/color"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/validate"
)

func WriteStampedPDF(ctx *model.Context, outputPDF io.Writer) error {
	return api.WriteContext(ctx, outputPDF)
}

// BuildStampedContext // take inputPDF and  return *model.Context
func BuildStampedContext(textStamps []TextStamp, imageStamps []ImageStamp) (*model.Context, error) {
	// Ensure fonts are installed before any watermarking occurs
	if err := InstallFonts(); err != nil {
		return nil, err
	}

	template, err := Tax50tawiPDFTemplate()
	if err != nil {
		return nil, err
	}

	ctx, err := ReadContext(template)
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

// applyTextWatermark applies a text watermark with the given configuration
func applyTextWatermark(pdfCtx *model.Context, stamp TextStamp) error {
	wm, err := TextWatermark(stamp)
	if err != nil {
		return err
	}
	return api.WatermarkContext(pdfCtx, nil, wm)
}

func TextWatermark(stamp TextStamp) (*model.Watermark, error) {
	wm, err := pdfcpu.ParseTextWatermarkDetails(ifEmpty(stamp.Text), "", true, 1)
	if err != nil {
		return nil, err
	}

	font := "THSarabunNew"
	if stamp.FontName != "" {
		font = stamp.FontName
	}

	wm.FillColor = color.Black
	wm.Dy = stamp.Dy
	wm.Dx = stamp.Dx
	wm.Diagonal = 0
	wm.Rotation = 0
	wm.Scale = 1
	wm.ScaleAbs = true
	wm.FontName = font
	wm.FontSize = stamp.FontSize
	wm.OnTop = true
	wm.Pos = stamp.Position

	return wm, nil
}

func applyImageWatermark(pdfCtx *model.Context, stamp ImageStamp) error {
	wm, err := ImageWatermark(stamp)
	if err != nil {
		return err
	}

	return api.WatermarkContext(pdfCtx, nil, wm)
}

// Stamp Image take reader
func ImageWatermark(stamp ImageStamp) (*model.Watermark, error) {
	wm, err := api.ImageWatermarkForReader(stamp.Reader, "", true, false, types.POINTS)
	if err != nil {
		return nil, err
	}
	wm.Dy = stamp.Dy
	wm.Dx = stamp.Dx
	wm.Scale = stamp.Scale
	wm.ScaleAbs = true
	wm.Opacity = stamp.Opacity
	wm.Diagonal = stamp.Diagonal
	wm.Rotation = 0
	wm.OnTop = stamp.OnTop
	wm.Pos = stamp.Pos
	return wm, nil
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
