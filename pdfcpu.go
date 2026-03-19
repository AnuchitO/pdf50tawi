package pdf50tawi

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	"github.com/signintech/gopdf"
)

const (
	pageWidth  = 595.28
	pageHeight = 841.89
)

// anchorToXY converts a pdfcpu-style anchor+offset to gopdf absolute coordinates.
// pdfcpu uses PDF standard axes (y up from bottom-left).
// gopdf uses screen axes (y down from top-left).
func anchorToXY(anchor Anchor, dx, dy float64) (float64, float64) {
	switch anchor {
	case TopLeft:
		return dx, -dy
	case TopCenter:
		return pageWidth/2 + dx, -dy
	case TopRight:
		return pageWidth + dx, -dy
	case BottomLeft:
		return dx, pageHeight - dy
	case BottomCenter:
		return pageWidth/2 + dx, pageHeight - dy
	case BottomRight:
		return pageWidth + dx, pageHeight - dy
	case Center:
		return pageWidth/2 + dx, pageHeight/2 - dy
	default:
		return dx, dy
	}
}

// stampPDF builds the output PDF by importing the template, then placing all
// text and image stamps. The Thai font is embedded once with subsetting.
func stampPDF(textStamps []TextStamp, imageStamps []ImageStamp, out io.Writer) error {
	tpl, err := Tax50tawiPDFTemplate()
	if err != nil {
		return err
	}
	tplData, err := io.ReadAll(tpl)
	if err != nil {
		return fmt.Errorf("read template: %w", err)
	}
	tplFile, err := os.CreateTemp("", "pdf50tawi-tpl-*.pdf")
	if err != nil {
		return fmt.Errorf("create temp: %w", err)
	}
	tplPath := tplFile.Name()
	defer os.Remove(tplPath)
	if _, err := tplFile.Write(tplData); err != nil {
		tplFile.Close()
		return fmt.Errorf("write temp template: %w", err)
	}
	tplFile.Close()

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: pageWidth, H: pageHeight}})

	if err := pdf.AddTTFFontData("THSarabunNew", thSarabunFontData); err != nil {
		return fmt.Errorf("add Thai font: %w", err)
	}

	tplIdx := pdf.ImportPage(tplPath, 1, "/MediaBox")
	pdf.AddPage()
	pdf.UseImportedTemplate(tplIdx, 0, 0, pageWidth, pageHeight)

	for i, img := range imageStamps {
		if err := placeImage(&pdf, img, i); err != nil {
			return err
		}
	}

	for _, stamp := range textStamps {
		if err := placeText(&pdf, stamp); err != nil {
			return err
		}
	}

	_, err = pdf.WriteTo(out)
	return err
}

func placeText(pdf *gopdf.GoPdf, stamp TextStamp) error {
	if err := pdf.SetFont("THSarabunNew", "", float64(stamp.FontSize)); err != nil {
		return fmt.Errorf("set font: %w", err)
	}
	x, y := anchorToXY(stamp.Position, stamp.Dx, stamp.Dy)
	pdf.SetXY(x, y)
	return pdf.Text(stamp.Text)
}

func placeImage(pdf *gopdf.GoPdf, stamp ImageStamp, idx int) error {
	if stamp.Reader == nil {
		return nil
	}
	data, err := io.ReadAll(stamp.Reader)
	if err != nil {
		return err
	}

	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || cfg.Width == 0 {
		return nil // skip invalid/empty images
	}

	w := pageWidth * stamp.Scale
	h := w * float64(cfg.Height) / float64(cfg.Width)
	x, y := anchorToXY(stamp.Pos, stamp.Dx, stamp.Dy)

	ext := ".png"
	if len(data) > 2 && data[0] == 0xFF && data[1] == 0xD8 {
		ext = ".jpg"
	}

	tmp, err := os.CreateTemp("", fmt.Sprintf("pdf-img%d-*%s", idx, ext))
	if err != nil {
		return err
	}
	imgPath := tmp.Name()
	defer os.Remove(imgPath)

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	tmp.Close()

	return pdf.Image(imgPath, x, y, &gopdf.Rect{W: w, H: h})
}
