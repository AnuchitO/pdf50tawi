package pdf50tawi

import "testing"

func TestInstallFonts(t *testing.T) {
	if err := InstallFonts(); err != nil {
		t.Fatalf("InstallFonts error: %v", err)
	}
}
