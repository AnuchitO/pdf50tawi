package pdf50tawi

import (
	"embed"
	"fmt"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/font"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

//go:embed fonts
var fontFS embed.FS

func InstallFonts() error {
	// Initialize pdfcpu configuration to set up UserFontDir
	if err := model.EnsureDefaultConfigAt(os.TempDir(), false); err != nil {
		return fmt.Errorf("failed to initialize pdfcpu config: %v", err)
	}

	// Install all THSarabunNew font variants to ensure complete glyph coverage
	fts := []string{"THSarabunNew.ttf", "THSarabunNew Bold.ttf", "THSarabunNew Italic.ttf", "THSarabunNew BoldItalic.ttf"}
	for _, ft := range fts {
		fontData, err := fontFS.ReadFile("fonts/THSarabunNew/" + ft)
		if err != nil {
			return fmt.Errorf("failed to read font %s: %v", ft, err)
		}

		// InstallFontFromBytes processes the complete TrueType font file
		// including all glyphs, character mappings, and Unicode support
		if err := font.InstallFontFromBytes(font.UserFontDir, ft, fontData); err != nil {
			return fmt.Errorf("failed to install font %s: %v", ft, err)
		}
		fmt.Printf("Installed font: %s\n", ft)
	}

	// Load user fonts to make them available for use
	// This step is crucial for glyph lookup and text rendering
	if err := font.LoadUserFonts(); err != nil {
		return fmt.Errorf("failed to load user fonts: %v", err)
	}

	return nil
}
