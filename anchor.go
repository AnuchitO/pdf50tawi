package pdf50tawi

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
