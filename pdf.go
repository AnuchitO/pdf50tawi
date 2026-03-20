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


// drawCheckmark draws a professional-quality filled ✓ using variable-width stroke expansion.
//
// Design (font-designer approach):
//   - Center line: straight left arm + cubic Bézier right arm (smooth sweep).
//   - Stroke width tapers from hMin at tips to hMax at valley (calligraphic weight).
//   - Outline = N sampled perpendicular offsets → dense polygon → smooth curves, no artifacts.
//
// Coordinate space: gopdf screen (y↓).
func drawCheckmark(pdf *gopdf.GoPdf, x, y, size float64) error {
	const N = 16             // samples per arm; higher = smoother
	hMax := size * 0.115    // half-width at valley (thickest point)
	hMin := size * 0.028    // half-width at tips   (tapered to near-zero)

	// ── Center-line key points ──────────────────────────────────────────────
	// Left arm:  straight  P0 → P1
	p0x, p0y := x+0.10*size, y+0.50*size // left tip
	p1x, p1y := x+0.32*size, y+0.82*size // valley

	// Right arm: cubic Bézier  P1 → C1 → C2 → P2
	// C1 directly right of valley so the arm starts horizontally (smooth join).
	c1x, c1y := x+0.46*size, y+0.82*size // ctrl-1: horizontal start
	c2x, c2y := x+0.76*size, y+0.22*size // ctrl-2: pulls toward tip
	p2x, p2y := x+0.96*size, y+0.07*size // right tip

	// ── Helpers ──────────────────────────────────────────────────────────────
	// For tangent (dx,dy), CCW perp = outer boundary, CW perp = inner boundary.
	outer := make([]gopdf.Point, 0, N*2+4)
	inner := make([]gopdf.Point, 0, N*2+4)

	addSample := func(px, py, dx, dy, h float64) {
		// normalize tangent
		l := math.Sqrt(dx*dx + dy*dy)
		if l < 1e-9 {
			return
		}
		dx /= l
		dy /= l
		// outer = CCW perp (-dy, dx); inner = CW perp (dy, -dx)
		outer = append(outer, gopdf.Point{X: px - dy*h, Y: py + dx*h})
		inner = append(inner, gopdf.Point{X: px + dy*h, Y: py - dx*h})
	}

	// ── Left arm (straight, t: 0=tip → 1=valley) ────────────────────────────
	ldx, ldy := p1x-p0x, p1y-p0y
	for i := 0; i <= N; i++ {
		t := float64(i) / float64(N)
		h := hMin + t*(hMax-hMin)
		addSample(p0x+t*ldx, p0y+t*ldy, ldx, ldy, h)
	}

	// ── Right arm (Bézier, t: 0=valley → 1=tip) ─────────────────────────────
	for i := 1; i <= N; i++ {
		t := float64(i) / float64(N)
		u := 1 - t
		// point on cubic Bézier
		px := u*u*u*p1x + 3*u*u*t*c1x + 3*u*t*t*c2x + t*t*t*p2x
		py := u*u*u*p1y + 3*u*u*t*c1y + 3*u*t*t*c2y + t*t*t*p2y
		// tangent (derivative)
		tdx := 3 * (u*u*(c1x-p1x) + 2*u*t*(c2x-c1x) + t*t*(p2x-c2x))
		tdy := 3 * (u*u*(c1y-p1y) + 2*u*t*(c2y-c1y) + t*t*(p2y-c2y))
		h := hMax + t*(hMin-hMax) // thick→thin
		addSample(px, py, tdx, tdy, h)
	}

	// ── Build polygon: outer forward + inner reversed ─────────────────────────
	pts := make([]gopdf.Point, 0, len(outer)+len(inner))
	pts = append(pts, outer...)
	for i := len(inner) - 1; i >= 0; i-- {
		pts = append(pts, inner[i])
	}

	pdf.SetFillColor(0, 0, 0)
	pdf.Polygon(pts, "F")
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
