package stamp

import (
	"io"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/color"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/draw"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// TextStampConfig represents the configuration for a text stamp.
type TextStampConfig struct {
	// Font configuration
	FontName string // Font family (e.g., "THSarabunNew-Bold")
	FontSize int    // Font size in points (e.g., 14)

	// Position and offset
	Position types.Anchor // Position anchor (e.g., tl, tr, bl, br, c, etc.)
	XOffset  float64      // Horizontal offset from position
	YOffset  float64      // Vertical offset from position

	// Appearance
	Scale      float64           // Scale factor (e.g., 0.2)
	Rotation   float64           // Rotation in degrees
	FillColor  color.SimpleColor // Text fill color
	RenderMode draw.RenderMode   // Text rendering mode (fill, stroke, etc.)

	// Additional watermark configuration
	Opacity float64 // Opacity (0.0 to 1.0)
	OnTop   bool    // If true, renders as stamp (on top), otherwise as watermark
}

// DefaultTextStampConfig returns a TextStampConfig with default values.
func DefaultTextStampConfig() TextStampConfig {
	return TextStampConfig{
		FontName:   "THSarabunNew",
		FontSize:   14,
		Position:   types.Center,
		Scale:      1.0,
		FillColor:  color.Black,
		RenderMode: draw.RMFill,
		Opacity:    1.0,
		OnTop:      true, // Default to stamp (on top)
	}
}

// ToWatermark converts TextStampConfig to a model.Watermark
func (c TextStampConfig) ToWatermark(text string) *model.Watermark {
	wm := model.DefaultWatermarkConfig()
	wm.OnTop = c.OnTop
	wm.TextString = text
	wm.FontName = c.FontName
	wm.FontSize = c.FontSize
	wm.Pos = c.Position
	wm.Dx = c.XOffset
	wm.Dy = c.YOffset
	wm.Scale = c.Scale
	wm.ScaleAbs = true
	wm.Rotation = c.Rotation
	wm.FillColor = c.FillColor
	wm.RenderMode = c.RenderMode
	wm.Opacity = c.Opacity
	return wm
}

// ParseTextStampConfig parses a text stamp configuration string into a TextStampConfig.
// Format: "font:<name>, points:<size>, scale:<scale>, pos:<position>, off:<x> <y>, rot:<degrees>, fill:<color>"
func ParseTextStampConfig(configStr string) (TextStampConfig, error) {
	cfg := DefaultTextStampConfig()
	// TODO: Implement parsing logic for the config string
	return cfg, nil
}

// ImageStampConfig represents the configuration for an image stamp.
type ImageStampConfig struct {
	// Image source
	FileName string    // Path to the image file  TODO: maybe remove this
	Image    io.Reader // Image reader (alternative to FileName) TODO: maybe remove this

	// Position and offset
	Position types.Anchor // Position anchor (e.g., tl, tr, bl, br, c, etc.)
	XOffset  float64      // Horizontal offset from position
	YOffset  float64      // Vertical offset from position

	// Appearance
	Scale    float64 // Scale factor (0.0 to 1.0 for relative, >1 for absolute)
	ScaleAbs bool    // If true, scale is absolute
	Rotation float64 // Rotation in degrees
	Opacity  float64 // Opacity (0.0 to 1.0)
	OnTop    bool    // If true, renders as stamp (on top), otherwise as watermark

	// Image dimensions (will be set automatically if not provided)
	Width  int
	Height int
}

// DefaultImageStampConfig returns an ImageStampConfig with default values.
func DefaultImageStampConfig() ImageStampConfig {
	return ImageStampConfig{
		Position: types.Center,
		Scale:    0.5, // Default to 50% of page dimension
		Opacity:  1.0,
		OnTop:    true,
	}
}

// ToWatermark converts ImageStampConfig to a model.Watermark
func (c ImageStampConfig) ToWatermark() *model.Watermark {
	wm := model.DefaultWatermarkConfig()
	wm.OnTop = c.OnTop
	wm.Mode = model.WMImage
	wm.FileName = c.FileName
	wm.Image = c.Image
	wm.Pos = c.Position
	wm.Dx = c.XOffset
	wm.Dy = c.YOffset
	wm.Scale = c.Scale
	wm.ScaleAbs = c.ScaleAbs
	wm.Rotation = c.Rotation
	wm.Opacity = c.Opacity

	// Set dimensions if provided
	if c.Width > 0 && c.Height > 0 {
		wm.Width = c.Width
		wm.Height = c.Height
	}

	return wm
}

// ParseImageStampConfig parses an image stamp configuration string into an ImageStampConfig.
// Format: "file:<path>, scale:<scale>, pos:<position>, off:<x> <y>, rot:<degrees>, opacity:<value>"
func ParseImageStampConfig(configStr string) (ImageStampConfig, error) {
	cfg := DefaultImageStampConfig()
	// TODO: Implement parsing logic for the config string
	return cfg, nil
}
