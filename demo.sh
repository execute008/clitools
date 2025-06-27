#!/bin/bash

# Demo script for clitools image optimization

echo "Building clitools..."
go build -o clitools .

echo ""
echo "Creating test images..."

# Create test image with Go
cat > create_demo_image.go << 'EOF'
package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	// Create a test image with significant transparent padding
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	
	// Fill with transparent background
	for y := 0; y < 200; y++ {
		for x := 0; x < 200; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 0})
		}
	}
	
	// Draw a colorful shape in the center with lots of transparent padding
	for y := 50; y < 150; y++ {
		for x := 50; x < 150; x++ {
			// Create a gradient effect
			r := uint8((x - 50) * 255 / 100)
			g := uint8((y - 50) * 255 / 100)
			b := uint8(((x - 50) + (y - 50)) * 255 / 200)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
	
	// Save as PNG
	file, err := os.Create("demo-image.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	
	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
	
	println("Created demo-image.png")
}
EOF

go run create_demo_image.go
rm create_demo_image.go

echo ""
echo "Original image size:"
ls -lh demo-image.png

echo ""
echo "Optimizing image with clitools..."
./clitools image optimize demo-image.png demo-optimized.webp -q 85

echo ""
echo "Comparing file sizes:"
echo "Original PNG: $(ls -lh demo-image.png | awk '{print $5}')"
echo "Optimized WebP: $(ls -lh demo-optimized.webp | awk '{print $5}')"

echo ""
echo "Demo complete! Files created:"
echo "- demo-image.png (original with transparent padding)"
echo "- demo-optimized.webp (cropped and converted to WebP)"
echo ""
echo "You can now test with your own images:"
echo "./clitools image optimize your-image.png output.webp"
