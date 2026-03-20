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
	const N = 28  // samples per Bézier arm (higher = smoother)
	const C = 10  // semicircle cap divisions (C+1 points per cap, 180/C° per step)
	h := size * 0.165 // half stroke width

	// ── Center line ─────────────────────────────────────────────────────────
	p0x, p0y := x+0.06*size, y+0.30*size // left tip
	p1x, p1y := x+0.24*size, y+0.86*size // valley
	c1x, c1y := x+0.42*size, y+0.92*size // Bézier ctrl1: pull down for wide open valley
	c2x, c2y := x+0.80*size, y+0.26*size // Bézier ctrl2
	p2x, p2y := x+0.97*size, y+0.08*size // right tip

	// ── Left arm: direction + perps ──────────────────────────────────────────
	ldx, ldy := p1x-p0x, p1y-p0y
	ll := math.Sqrt(ldx*ldx + ldy*ldy); ldx /= ll; ldy /= ll
	lox, loy := -ldy, ldx  // outer (CCW of forward)
	lix, liy := ldy, -ldx  // inner (CW of forward)
	lbx, lby := -ldx, -ldy // backward (true tip direction for left cap)

	// ── Right arm: tangent at tip (t=1) ──────────────────────────────────────
	rtdx, rtdy := p2x-c2x, p2y-c2y
	rl := math.Sqrt(rtdx*rtdx + rtdy*rtdy); rtdx /= rl; rtdy /= rl
	rox, roy := -rtdy, rtdx // outer (CCW of forward); cap ends at inner = -(rox,roy)

	type pt = gopdf.Point
	poly := make([]pt, 0, 2*N+2*C+30)

	// bezOI: sample cubic Bézier (p1→c1→c2→p2) at t, return outer/inner offset pts
	bezOI := func(t, hh float64) (po, pi pt) {
		u := 1 - t
		px := u*u*u*p1x + 3*u*u*t*c1x + 3*u*t*t*c2x + t*t*t*p2x
		py := u*u*u*p1y + 3*u*u*t*c1y + 3*u*t*t*c2y + t*t*t*p2y
		tdx := 3 * (u*u*(c1x-p1x) + 2*u*t*(c2x-c1x) + t*t*(p2x-c2x))
		tdy := 3 * (u*u*(c1y-p1y) + 2*u*t*(c2y-c1y) + t*t*(p2y-c2y))
		tl := math.Sqrt(tdx*tdx + tdy*tdy); tdx /= tl; tdy /= tl
		return pt{px - tdy*hh, py + tdx*hh}, pt{px + tdy*hh, py - tdx*hh}
	}

	// smoothCap: C+1 points sweeping a semicircle from direction a1 through fwd to -a1.
	// cos(θ)*a1 + sin(θ)*fwd for θ in [0, π].
	smoothCap := func(cx, cy, a1x, a1y, fwdx, fwdy float64) {
		for i := 0; i <= C; i++ {
			θ := math.Pi * float64(i) / float64(C)
			c, s := math.Cos(θ), math.Sin(θ)
			poly = append(poly, pt{cx + (c*a1x+s*fwdx)*h, cy + (c*a1y+s*fwdy)*h})
		}
	}

	// ── 1. Left arm outer: tip → valley ──────────────────────────────────────
	for i := 0; i <= N; i++ {
		t := float64(i) / float64(N)
		poly = append(poly, pt{p0x + t*(p1x-p0x) + lox*h, p0y + t*(p1y-p0y) + loy*h})
	}

	// ── 2. Valley outer apex ──────────────────────────────────────────────────
	poly = append(poly, pt{p1x, p1y + h})

	// ── 3. Right arm outer: valley → tip ─────────────────────────────────────
	for i := 1; i <= N; i++ {
		o, _ := bezOI(float64(i)/float64(N), h)
		poly = append(poly, o)
	}

	// ── 4. Right cap: sweep rox → rtdx → rix (forward = true tip end) ────────
	smoothCap(p2x, p2y, rox, roy, rtdx, rtdy)

	// ── 5. Right arm inner: tip → valley ─────────────────────────────────────
	for i := N - 1; i >= 1; i-- {
		_, pi := bezOI(float64(i)/float64(N), h)
		poly = append(poly, pi)
	}

	// ── 6. Valley inner: wide open quadratic Bézier arc ──────────────────────
	rivX, rivY := p1x, p1y-h                  // right-arm inner at valley
	livX, livY := p1x+lix*h, p1y+liy*h        // left-arm inner at valley
	ctX := (rivX + livX) * 0.5
	ctY := p1y + h*1.2 // control pulled below valley → wide flat U
	for i := 0; i <= 8; i++ {
		t := float64(i) / 8.0
		u := 1 - t
		poly = append(poly, pt{
			u*u*rivX + 2*u*t*ctX + t*t*livX,
			u*u*rivY + 2*u*t*ctY + t*t*livY,
		})
	}

	// ── 7. Left arm inner: valley → tip ──────────────────────────────────────
	for i := N - 1; i >= 0; i-- {
		t := float64(i) / float64(N)
		poly = append(poly, pt{p0x + t*(p1x-p0x) + lix*h, p0y + t*(p1y-p0y) + liy*h})
	}

	// ── 8. Left cap: sweep lix → lbx → lox (backward = true tip end) ─────────
	smoothCap(p0x, p0y, lix, liy, lbx, lby)

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
