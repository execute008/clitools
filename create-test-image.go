package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	// Create a test image with transparent padding
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	
	// Fill with transparent background
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 0})
		}
	}
	
	// Draw a small red square in the center (leaving transparent padding)
	for y := 30; y < 70; y++ {
		for x := 30; x < 70; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}
	
	// Save as PNG
	file, err := os.Create("test-image.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	
	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
	
	println("Created test-image.png with transparent padding")
}
