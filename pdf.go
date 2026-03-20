package pdf50tawi

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
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


// drawCheckmark draws a bold ✓ matching the reference style:
// thick uniform stroke, large rounded caps at both tips, smooth rounded valley.
//
// The full outline is built in one clockwise pass:
//
//	left-arm-outer → valley-outer → right-arm-outer →
//	right-cap → right-arm-inner → valley-inner-arc →
//	left-arm-inner → left-cap → (close)
func drawCheckmark(pdf *gopdf.GoPdf, x, y, size float64) error {
	const N = 14           // samples per arm
	h := size * 0.155      // half stroke width (bold, uniform)

	// ── Center line ─────────────────────────────────────────────────────────
	p0x, p0y := x+0.12*size, y+0.44*size // left tip
	p1x, p1y := x+0.32*size, y+0.80*size // valley
	c1x, c1y := x+0.47*size, y+0.80*size // Bézier ctrl1: horizontal departure
	c2x, c2y := x+0.74*size, y+0.18*size // Bézier ctrl2
	p2x, p2y := x+0.95*size, y+0.06*size // right tip

	// ── Left arm unit direction + perps ─────────────────────────────────────
	ldx, ldy := p1x-p0x, p1y-p0y
	ll := math.Sqrt(ldx*ldx + ldy*ldy); ldx /= ll; ldy /= ll
	lox, loy := -ldy, ldx  // outer (CCW)
	lix, liy := ldy, -ldx  // inner (CW)
	lbx, lby := -ldx, -ldy // backward (for left cap)

	// ── Right arm tangent at tip (t=1) ──────────────────────────────────────
	rtdx, rtdy := p2x-c2x, p2y-c2y // ∝ B'(1)
	rl := math.Sqrt(rtdx*rtdx + rtdy*rtdy); rtdx /= rl; rtdy /= rl
	rox, roy := -rtdy, rtdx  // outer at tip
	rix, riy := rtdy, -rtdx  // inner at tip
	rbx, rby := -rtdx, -rtdy // backward at tip

	type pt = gopdf.Point
	poly := make([]pt, 0, 4*N+30)

	// helper: sample Bézier point + CCW/CW offsets
	bezierOuterInner := func(t, hh float64) (po, pi pt) {
		u := 1 - t
		px := u*u*u*p1x + 3*u*u*t*c1x + 3*u*t*t*c2x + t*t*t*p2x
		py := u*u*u*p1y + 3*u*u*t*c1y + 3*u*t*t*c2y + t*t*t*p2y
		tdx := 3 * (u*u*(c1x-p1x) + 2*u*t*(c2x-c1x) + t*t*(p2x-c2x))
		tdy := 3 * (u*u*(c1y-p1y) + 2*u*t*(c2y-c1y) + t*t*(p2y-c2y))
		tl := math.Sqrt(tdx*tdx + tdy*tdy); tdx /= tl; tdy /= tl
		return pt{px - tdy*hh, py + tdx*hh}, pt{px + tdy*hh, py - tdx*hh}
	}

	// cap45 returns the 45° interpolated cap point (unit vectors must be perpendicular)
	cap45 := func(cx, cy, ax, ay, bx, by, hh float64) pt {
		return pt{cx + (ax+bx)*hh/math.Sqrt2, cy + (ay+by)*hh/math.Sqrt2}
	}

	// ── 1. Left arm outer: tip → valley ─────────────────────────────────────
	for i := 0; i <= N; i++ {
		t := float64(i) / float64(N)
		poly = append(poly, pt{p0x + t*(p1x-p0x) + lox*h, p0y + t*(p1y-p0y) + loy*h})
	}

	// ── 2. Valley outer: right-arm starts horizontally, outer = (0,+1) ──────
	poly = append(poly, pt{p1x, p1y + h})

	// ── 3. Right arm outer: valley → tip ────────────────────────────────────
	for i := 1; i <= N; i++ {
		o, _ := bezierOuterInner(float64(i)/float64(N), h)
		poly = append(poly, o)
	}

	// ── 4. Right cap (5-point semicircle) ────────────────────────────────────
	poly = append(poly,
		pt{p2x + rox*h, p2y + roy*h},
		cap45(p2x, p2y, rox, roy, rbx, rby, h),
		pt{p2x + rbx*h, p2y + rby*h},
		cap45(p2x, p2y, rbx, rby, rix, riy, h),
		pt{p2x + rix*h, p2y + riy*h},
	)

	// ── 5. Right arm inner: tip → valley ────────────────────────────────────
	for i := N - 1; i >= 1; i-- {
		_, pi := bezierOuterInner(float64(i)/float64(N), h)
		poly = append(poly, pi)
	}

	// ── 6. Valley inner: smooth rounded arc ──────────────────────────────────
	// Right arm inner at valley (tangent is horizontal → CW inner = (0,−1))
	rivX, rivY := p1x, p1y-h
	// Left arm inner at valley
	livX, livY := p1x+lix*h, p1y+liy*h
	// Quadratic Bézier control pulled toward valley center → creates the concave arc
	ctX := (rivX+livX)*0.45 + p1x*0.55
	ctY := (rivY+livY)*0.45 + p1y*0.55
	for _, t := range []float64{0.0, 0.25, 0.5, 0.75, 1.0} {
		u := 1 - t
		poly = append(poly, pt{
			u*u*rivX + 2*u*t*ctX + t*t*livX,
			u*u*rivY + 2*u*t*ctY + t*t*livY,
		})
	}

	// ── 7. Left arm inner: valley → tip ─────────────────────────────────────
	for i := N - 1; i >= 0; i-- {
		t := float64(i) / float64(N)
		poly = append(poly, pt{p0x + t*(p1x-p0x) + lix*h, p0y + t*(p1y-p0y) + liy*h})
	}

	// ── 8. Left cap (5-point semicircle) ────────────────────────────────────
	poly = append(poly,
		cap45(p0x, p0y, lix, liy, lbx, lby, h),
		pt{p0x + lbx*h, p0y + lby*h},
		cap45(p0x, p0y, lbx, lby, lox, loy, h),
	)

	pdf.SetFillColor(0, 0, 0)
	pdf.Polygon(poly, "F")
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
