package tray

import (
	"bytes"
	"image/png"
	"testing"
)

func TestDefaultIconBytes(t *testing.T) {
	iconBytes := DefaultIconBytes()
	if len(iconBytes) == 0 {
		t.Fatal("expected non-empty icon bytes")
	}

	// Verify valid PNG
	img, err := png.Decode(bytes.NewReader(iconBytes))
	if err != nil {
		t.Fatalf("failed to decode generated PNG icon: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != 32 || bounds.Dy() != 32 {
		t.Fatalf("expected 32x32 icon, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}
