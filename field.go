package pdf50tawi

import "io"

// Anchor represents a reference point on a PDF page for positioning text and images.
type Anchor int

const (
	TopLeft     Anchor = iota
	TopCenter          // 1
	TopRight           // 2
	Left               // 3
	Center             // 4
	Right              // 5
	BottomLeft         // 6
	BottomCenter       // 7
	BottomRight        // 8
)

// TextField defines a text value and its position on the certificate form.
type TextField struct {
	Text     string
	Dx       float64
	Dy       float64
	FontSize int
	FontName string
	Position Anchor
}

// ImageField defines an image (signature or seal) and its position on the certificate form.
type ImageField struct {
	Reader   io.Reader
	Pos      Anchor
	Dx       float64
	Dy       float64
	Scale    float64
	Opacity  float64
	Diagonal int
	OnTop    bool
}
