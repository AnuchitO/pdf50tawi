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
	// Register a system font for checkmarks (✓). THSarabunNew doesn't have U+2713.

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
	x, y := anchorToXY(stamp.Position, stamp.Dx, stamp.Dy)

	// ✓ has no glyph in THSarabunNew — draw as a filled vector polygon instead.
	if stamp.Text == "✓" {
		return drawCheckmark(pdf, x, y, float64(stamp.FontSize))
	}

	if err := pdf.SetFont("THSarabunNew", "", float64(stamp.FontSize)); err != nil {
		return fmt.Errorf("set font: %w", err)
	}

	// pdfcpu anchors the text bounding box corner that matches the anchor name.
	// gopdf.Text() always starts text at the left edge, so we must shift x to
	// replicate right-align (BottomRight/TopRight/Right) and center (XCenter).
	switch stamp.Position {
	case TopCenter, BottomCenter, Center:
		if w, err := pdf.MeasureTextWidth(stamp.Text); err == nil {
			x -= w / 2
		}
	case TopRight, BottomRight, Right:
		if w, err := pdf.MeasureTextWidth(stamp.Text); err == nil {
			x -= w
		}
	}

	pdf.SetXY(x, y)
	return pdf.Text(stamp.Text)
}


// drawCheckmark draws a filled ✓ at (x,y) with rounded caps at both tips.
// The shape is a 12-point polygon: 5 points per rounded cap + 1 point per valley side.
//
// Geometry (all coords in gopdf screen space, y↓):
//
//	P1 = left arm tip  (0.12s, 0.48s)
//	P2 = valley        (0.33s, 0.85s)
//	P3 = right arm tip (0.98s, 0.05s)
//
// Left arm dir  ≈ (0.560, 0.829)  inner-perp=(0.829,−0.560)  outer-perp=(−0.829, 0.560)
// Right arm dir ≈ (0.631,−0.776)  inner-perp=(−0.776,−0.631) outer-perp=( 0.776, 0.631)
// Valley bisector outer ≈ (0,+1), inner ≈ (0,−1).
func drawCheckmark(pdf *gopdf.GoPdf, x, y, size float64) error {
	h := size * 0.12

	p1x, p1y := x+0.12*size, y+0.48*size
	p2x, p2y := x+0.33*size, y+0.85*size
	p3x, p3y := x+0.98*size, y+0.05*size

	points := []gopdf.Point{
		// — Left arm rounded cap (5 pts, inner → tip → outer) —
		{X: p1x + 0.829*h, Y: p1y - 0.560*h}, // inner edge
		{X: p1x + 0.190*h, Y: p1y - 0.982*h}, // 45° toward tip
		{X: p1x - 0.560*h, Y: p1y - 0.829*h}, // cap tip
		{X: p1x - 0.982*h, Y: p1y - 0.190*h}, // 45° toward outer
		{X: p1x - 0.829*h, Y: p1y + 0.560*h}, // outer edge
		// — Valley outer (bottom of valley) —
		{X: p2x, Y: p2y + h},
		// — Right arm rounded cap (5 pts, outer → tip → inner) —
		{X: p3x + 0.776*h, Y: p3y + 0.631*h}, // outer edge
		{X: p3x + 0.994*h, Y: p3y - 0.103*h}, // 45° toward tip
		{X: p3x + 0.631*h, Y: p3y - 0.776*h}, // cap tip
		{X: p3x - 0.103*h, Y: p3y - 0.994*h}, // 45° toward inner
		{X: p3x - 0.776*h, Y: p3y - 0.631*h}, // inner edge
		// — Valley inner (top of valley) —
		{X: p2x, Y: p2y - h},
	}

	pdf.SetFillColor(0, 0, 0)
	pdf.Polygon(points, "F")
	return nil
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
