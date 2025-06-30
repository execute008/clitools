#!/bin/bash

# Demo script for clitools image optimization

echo ""
echo ""
echo ""
echo "=========================================="
echo "    CLI Tools - Image Optimization Demo  "
echo "=========================================="
echo ""
echo ""

echo "Building clitools..."
go build -o clitools .
echo "✅ Build complete!"
echo ""
echo ""
echo ""
echo "Press Enter to continue..."
read -r

echo ""
echo ""
echo ""
echo "=========================================="
echo "         Creating Test Images             "
echo "=========================================="
echo ""
echo ""

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

echo "✅ Test PNG created with transparent padding"
echo ""
echo ""
echo ""
echo "Press Enter to view file size..."
read -r

echo ""
echo ""
echo ""
echo "=========================================="
echo "         Original Image Analysis         "
echo "=========================================="
echo ""
echo ""
echo "Original image size:"
ls -lh demo-image.png
echo ""
echo "This PNG has lots of transparent padding that we can crop out!"
echo ""
echo ""
echo ""
echo "Press Enter to start optimization..."
read -r

echo ""
echo ""
echo ""
echo "=========================================="
echo "         Image Optimization Process       "
echo "=========================================="
echo ""
echo ""
echo "Optimizing image with clitools..."
echo "• Cropping transparent areas"
echo "• Converting to WebP format"
echo "• Quality setting: 85%"
echo ""
./clitools image optimize demo-image.png demo-optimized.webp -q 85
echo ""
echo "✅ Optimization complete!"
echo ""
echo ""
echo ""
echo "Press Enter to see results..."
read -r

echo ""
echo ""
echo ""
echo "=========================================="
echo "           Results Comparison             "
echo "=========================================="
echo ""
echo ""
echo "File size comparison:"
echo "📁 Original PNG:    $(ls -lh demo-image.png | awk '{print $5}')"
echo "🚀 Optimized WebP:  $(ls -lh demo-optimized.webp | awk '{print $5}')"
echo ""

# Calculate size reduction
original_size=$(stat -f%z demo-image.png)
optimized_size=$(stat -f%z demo-optimized.webp)
reduction=$((100 - (optimized_size * 100 / original_size)))
echo "💾 Size reduction: ${reduction}%"
echo ""
echo ""
echo ""
echo "Press Enter to continue..."
read -r

echo ""
echo ""
echo ""
echo "=========================================="
echo "         SVG Support Demo                "
echo "=========================================="
echo ""
echo ""
echo "Creating a test SVG with gradients..."

# Create test SVG
cat > demo-icon.svg << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<svg width="200" height="200" viewBox="0 0 200 200" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <linearGradient id="grad" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" style="stop-color:#ff6b6b;stop-opacity:1" />
      <stop offset="50%" style="stop-color:#4dabf7;stop-opacity:1" />
      <stop offset="100%" style="stop-color:#69db7c;stop-opacity:1" />
    </linearGradient>
  </defs>
  <circle cx="100" cy="100" r="80" fill="url(#grad)" stroke="#333" stroke-width="4"/>
  <text x="100" y="110" text-anchor="middle" fill="white" font-size="18" font-weight="bold">SVG</text>
</svg>
EOF

echo "✅ SVG created with gradients and text"
echo ""
echo ""
echo "Testing SVG optimization with different quality levels..."
echo ""

# Test different SVG scales
echo "🔄 Converting with 2x scale (default)..."
./clitools image optimize demo-icon.svg demo-icon-2x.webp --svg-scale 2 --quality 85

echo ""
echo "🔄 Converting with 4x scale (high quality)..."
./clitools image optimize demo-icon.svg demo-icon-4x.webp --svg-scale 4 --quality 95

echo ""
echo "SVG Results:"
echo "📁 Original SVG:     $(ls -lh demo-icon.svg | awk '{print $5}')"
echo "🚀 WebP (2x scale):  $(ls -lh demo-icon-2x.webp | awk '{print $5}')"
echo "💎 WebP (4x scale):  $(ls -lh demo-icon-4x.webp | awk '{print $5}')"
echo ""
echo ""
echo ""
echo "Press Enter to see final summary..."
read -r

echo ""
echo ""
echo ""
echo "=========================================="
echo "              Demo Complete!              "
echo "=========================================="
echo ""
echo ""
echo "Files created:"
echo "📄 demo-image.png        - Original PNG with transparent padding"
echo "🚀 demo-optimized.webp   - Cropped and optimized WebP"
echo "🎨 demo-icon.svg         - Test SVG with gradients"
echo "📱 demo-icon-2x.webp     - SVG converted to WebP (2x quality)"
echo "💎 demo-icon-4x.webp     - SVG converted to WebP (4x quality)"
echo ""
echo ""
echo "🎯 Key Benefits:"
echo "   • Automatic transparent area cropping"
echo "   • WebP conversion for better compression"
echo "   • SVG support with configurable quality"
echo "   • Significant file size reductions"
echo ""
echo ""
echo "📖 Usage Examples:"
echo ""
echo "Basic optimization:"
echo "  ./clitools image optimize input.png output.webp"
echo ""
echo "High quality SVG conversion:"
echo "  ./clitools image optimize icon.svg icon.webp --svg-scale 4 --quality 95"
echo ""
echo "Batch processing with custom quality:"
echo "  ./clitools image optimize photo.jpg photo.webp --quality 80"
echo ""
echo ""
echo "🧹 Cleanup: Run 'rm demo-*' to remove demo files"
echo ""
echo ""
echo "=========================================="
echo "       Thank you for trying CLI Tools!    "
echo "=========================================="
echo ""
echo ""
