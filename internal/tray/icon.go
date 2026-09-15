package tray

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

// DefaultIconBytes generates a crisp 32x32 PNG icon in memory with Dapodik blue theme
func DefaultIconBytes() []byte {
	const size = 32
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Fill background with transparent
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)

	// Outer rounded circle: Dapodik blue (#2196F3)
	blueColor := color.RGBA{R: 33, G: 150, B: 243, A: 255}
	whiteColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	centerX, centerY, radius := 16.0, 16.0, 14.0
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := float64(x) - centerX
			dy := float64(y) - centerY
			distSq := dx*dx + dy*dy
			if distSq <= radius*radius {
				img.Set(x, y, blueColor)
			}
		}
	}

	// Inner symbol: Clean horizontal bridge bars
	// Top bar
	for x := 8; x <= 23; x++ {
		for y := 10; y <= 12; y++ {
			img.Set(x, y, whiteColor)
		}
	}
	// Middle pillar left
	for x := 10; x <= 12; x++ {
		for y := 13; y <= 21; y++ {
			img.Set(x, y, whiteColor)
		}
	}
	// Middle pillar right
	for x := 19; x <= 21; x++ {
		for y := 13; y <= 21; y++ {
			img.Set(x, y, whiteColor)
		}
	}
	// Bottom base
	for x := 7; x <= 24; x++ {
		for y := 20; y <= 22; y++ {
			img.Set(x, y, whiteColor)
		}
	}

	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}
