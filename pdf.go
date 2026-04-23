package pdf50tawi

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"

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

// fillCertificate builds the output PDF by importing the template, then placing all
// text and image fields. The Thai font is embedded once with subsetting.
func fillCertificate(textFields []TextField, imageFields []ImageField, out io.Writer) error {
	tplPath, err := cachedTemplatePath()
	if err != nil {
		return err
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: pageWidth, H: pageHeight}})

	if err := pdf.AddTTFFontData("THSarabunNew", thSarabunFontData); err != nil {
		return fmt.Errorf("add Thai font: %w", err)
	}
	// Register a system font for checkmarks (✓). THSarabunNew doesn't have U+2713.

	tplIdx := pdf.ImportPage(tplPath, 1, "/MediaBox")
	pdf.AddPage()
	pdf.UseImportedTemplate(tplIdx, 0, 0, pageWidth, pageHeight)

	for _, img := range imageFields {
		if err := placeImage(&pdf, img); err != nil {
			return err
		}
	}

	for _, field := range textFields {
		if err := placeText(&pdf, field); err != nil {
			return err
		}
	}

	_, err = pdf.WriteTo(out)
	return err
}

func placeText(pdf *gopdf.GoPdf, field TextField) error {
	x, y := anchorToXY(field.Position, field.Dx, field.Dy)

	// ✓ has no glyph in THSarabunNew — draw as a filled vector polygon instead.
	if field.Text == "✓" {
		return drawCheckmark(pdf, x, y, float64(field.FontSize))
	}

	if err := pdf.SetFont("THSarabunNew", "", float64(field.FontSize)); err != nil {
		return fmt.Errorf("set font: %w", err)
	}

	// pdfcpu anchors the text bounding box corner that matches the anchor name.
	// gopdf.Text() always starts text at the left edge, so we must shift x to
	// replicate right-align (BottomRight/TopRight/Right) and center (XCenter).
	switch field.Position {
	case TopCenter, BottomCenter, Center:
		if w, err := pdf.MeasureTextWidth(field.Text); err == nil {
			x -= w / 2
		}
	case TopRight, BottomRight, Right:
		if w, err := pdf.MeasureTextWidth(field.Text); err == nil {
			x -= w
		}
	}

	pdf.SetXY(x, y)
	return pdf.Text(field.Text)
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
	const N = 28  // samples per arm
	const C = 16  // semicircle cap divisions (180/16 ≈ 11° per step = smooth)
	h := size * 0.125 // half stroke width

	// ── Center line ─────────────────────────────────────────────────────────
	// Geometry matched to reference: shorter left arm (tip is lower on page),
	// wide open valley, long sweeping right arm.
	p0x, p0y := x+0.10*size, y+0.42*size // left tip (shorter arm, starts lower)
	p1x, p1y := x+0.27*size, y+0.77*size // valley
	c1x, c1y := x+0.44*size, y+0.83*size // Bézier ctrl1: gentle pull for wide valley
	c2x, c2y := x+0.80*size, y+0.23*size // Bézier ctrl2: long sweep
	p2x, p2y := x+0.97*size, y+0.07*size // right tip

	// ── Left arm: direction + perps ──────────────────────────────────────────
	ldx, ldy := p1x-p0x, p1y-p0y
	ll := math.Sqrt(ldx*ldx + ldy*ldy); ldx /= ll; ldy /= ll
	lox, loy := -ldy, ldx  // outer (CCW)
	lix, liy := ldy, -ldx  // inner (CW)
	lbx, lby := -ldx, -ldy // backward (for left cap)

	// ── Right arm: tangent at tip ─────────────────────────────────────────────
	rtdx, rtdy := p2x-c2x, p2y-c2y
	rl := math.Sqrt(rtdx*rtdx + rtdy*rtdy); rtdx /= rl; rtdy /= rl
	rox, roy := -rtdy, rtdx // outer at right tip

	// ── Right arm: tangent at valley (t=0) for valley outer arc ─────────────
	rvtdx, rvtdy := c1x-p1x, c1y-p1y
	rvl := math.Sqrt(rvtdx*rvtdx + rvtdy*rvtdy); rvtdx /= rvl; rvtdy /= rvl

	type pt = gopdf.Point
	poly := make([]pt, 0, 2*N+2*C+40)

	// bezOI: sample cubic Bézier p1→c1→c2→p2, return outer/inner offset pts
	bezOI := func(t, hh float64) (po, pi pt) {
		u := 1 - t
		px := u*u*u*p1x + 3*u*u*t*c1x + 3*u*t*t*c2x + t*t*t*p2x
		py := u*u*u*p1y + 3*u*u*t*c1y + 3*u*t*t*c2y + t*t*t*p2y
		tdx := 3 * (u*u*(c1x-p1x) + 2*u*t*(c2x-c1x) + t*t*(p2x-c2x))
		tdy := 3 * (u*u*(c1y-p1y) + 2*u*t*(c2y-c1y) + t*t*(p2y-c2y))
		tl := math.Sqrt(tdx*tdx + tdy*tdy); tdx /= tl; tdy /= tl
		return pt{px - tdy*hh, py + tdx*hh}, pt{px + tdy*hh, py - tdx*hh}
	}

	// smoothCap: sweep semicircle cos(θ)*a1 + sin(θ)*fwd for θ ∈ [0, π]
	smoothCap := func(cx, cy, a1x, a1y, fwdx, fwdy float64) {
		for i := 0; i <= C; i++ {
			θ := math.Pi * float64(i) / float64(C)
			c, s := math.Cos(θ), math.Sin(θ)
			poly = append(poly, pt{cx + (c*a1x+s*fwdx)*h, cy + (c*a1y+s*fwdy)*h})
		}
	}

	// quadArc: sample quadratic Bézier with 8 segments
	quadArc := func(ax, ay, cx, cy, bx, by float64) {
		for i := 0; i <= 8; i++ {
			t := float64(i) / 8.0; u := 1 - t
			poly = append(poly, pt{u*u*ax + 2*u*t*cx + t*t*bx, u*u*ay + 2*u*t*cy + t*t*by})
		}
	}

	// ── 1. Left arm outer: tip → valley ──────────────────────────────────────
	for i := 0; i <= N; i++ {
		t := float64(i) / float64(N)
		poly = append(poly, pt{p0x + t*(p1x-p0x) + lox*h, p0y + t*(p1y-p0y) + loy*h})
	}

	// ── 2. Valley outer: single bottom point below valley center ─────────────
	poly = append(poly, pt{p1x, p1y + h})

	// ── 3. Right arm outer: valley → tip ─────────────────────────────────────
	for i := 1; i <= N; i++ {
		o, _ := bezOI(float64(i)/float64(N), h)
		poly = append(poly, o)
	}

	// ── 4. Right cap: sweep rox → rtdx → -rox ────────────────────────────────
	smoothCap(p2x, p2y, rox, roy, rtdx, rtdy)

	// ── 5. Right arm inner: tip → valley ─────────────────────────────────────
	for i := N - 1; i >= 1; i-- {
		_, pi := bezOI(float64(i)/float64(N), h)
		poly = append(poly, pi)
	}

	// ── 6. Valley inner: smooth concave quadratic Bézier arc ─────────────────
	// Right arm inner at valley: CW of right arm tangent at t=0
	rivX, rivY := p1x+rvtdy*h, p1y-rvtdx*h
	livX, livY := p1x+lix*h, p1y+liy*h // left arm inner at valley
	// Control pulled slightly toward valley center (below endpoints in screen)
	ctX := (rivX + livX) * 0.5
	ctY := p1y + h*0.6 // pushed below valley center → smooth wide concave arc
	quadArc(rivX, rivY, ctX, ctY, livX, livY)

	// ── 7. Left arm inner: valley → tip ──────────────────────────────────────
	for i := N - 1; i >= 0; i-- {
		t := float64(i) / float64(N)
		poly = append(poly, pt{p0x + t*(p1x-p0x) + lix*h, p0y + t*(p1y-p0y) + liy*h})
	}

	// ── 8. Left cap: sweep lix → lbx → lox ───────────────────────────────────
	smoothCap(p0x, p0y, lix, liy, lbx, lby)

	pdf.SetFillColor(0, 0, 0)
	pdf.Polygon(poly, "F")
	return nil
}

func placeImage(pdf *gopdf.GoPdf, field ImageField) error {
	if field.Reader == nil {
		return nil
	}
	data, err := io.ReadAll(field.Reader)
	if err != nil {
		return err
	}

	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || cfg.Width == 0 {
		return nil // skip invalid/empty images
	}

	w := pageWidth * field.Scale
	h := w * float64(cfg.Height) / float64(cfg.Width)
	x, y := anchorToXY(field.Pos, field.Dx, field.Dy)

	holder, err := gopdf.ImageHolderByBytes(data)
	if err != nil {
		return err
	}
	return pdf.ImageByHolder(holder, x, y, &gopdf.Rect{W: w, H: h})
}
