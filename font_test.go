package pdf50tawi

import "testing"

func TestFontDataEmbedded(t *testing.T) {
	if len(thSarabunFontData) == 0 {
		t.Fatal("THSarabunNew font data is not embedded")
	}
	// TTF files start with 0x00010000 or "OTTO" (OTF)
	if len(thSarabunFontData) < 4 {
		t.Fatal("font data too small to be a valid TTF")
	}
}
