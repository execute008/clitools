package image

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

type Processor struct{}

func NewProcessor() *Processor {
	return &Processor{}
}

// OptimizeImage crops transparent areas and converts to WebP format
func (p *Processor) OptimizeImage(inputPath, outputPath string, quality float32) error {
	var img image.Image
	var err error

	// Check if input is SVG
	if strings.ToLower(filepath.Ext(inputPath)) == ".svg" {
		img, err = p.loadSVG(inputPath)
		if err != nil {
			return fmt.Errorf("failed to load SVG: %w", err)
		}
	} else {
		// Open and decode regular image formats
		inputFile, err := os.Open(inputPath)
		if err != nil {
			return fmt.Errorf("failed to open input file: %w", err)
		}
		defer inputFile.Close()

		img, _, err = image.Decode(inputFile)
		if err != nil {
			return fmt.Errorf("failed to decode image: %w", err)
		}
	}

	// Crop transparent areas
	croppedImg := p.cropTransparentAreas(img)

	// Convert to WebP and save
	return p.saveAsWebP(croppedImg, outputPath, quality)
}

// OptimizeImageWithScale crops transparent areas and converts to WebP format with configurable SVG scaling
func (p *Processor) OptimizeImageWithScale(inputPath, outputPath string, quality, svgScale float32) error {
	var img image.Image
	var err error

	// Check if input is SVG
	if strings.ToLower(filepath.Ext(inputPath)) == ".svg" {
		img, err = p.loadSVGWithScale(inputPath, svgScale)
		if err != nil {
			return fmt.Errorf("failed to load SVG: %w", err)
		}
	} else {
		// Open and decode regular image formats
		inputFile, err := os.Open(inputPath)
		if err != nil {
			return fmt.Errorf("failed to open input file: %w", err)
		}
		defer inputFile.Close()

		img, _, err = image.Decode(inputFile)
		if err != nil {
			return fmt.Errorf("failed to decode image: %w", err)
		}
	}

	// Crop transparent areas
	croppedImg := p.cropTransparentAreas(img)

	// Convert to WebP and save
	return p.saveAsWebP(croppedImg, outputPath, quality)
}

// cropTransparentAreas removes transparent padding from the image
func (p *Processor) cropTransparentAreas(img image.Image) image.Image {
	bounds := img.Bounds()

	// Find the actual content bounds (non-transparent areas)
	minX, minY := bounds.Max.X, bounds.Max.Y
	maxX, maxY := bounds.Min.X, bounds.Min.Y

	foundContent := false

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Check if pixel is not transparent (alpha > 0) or has color content
			if a > 0 || r > 0 || g > 0 || b > 0 {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
				foundContent = true
			}
		}
	}

	// If no content found, return a 1x1 transparent image
	if !foundContent {
		return image.NewRGBA(image.Rect(0, 0, 1, 1))
	}

	// Create cropped image
	croppedBounds := image.Rect(0, 0, maxX-minX+1, maxY-minY+1)
	croppedImg := image.NewRGBA(croppedBounds)

	// Copy the non-transparent area to the new image
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			croppedImg.Set(x-minX, y-minY, img.At(x, y))
		}
	}

	return croppedImg
}

// saveAsWebP saves the image as WebP format
func (p *Processor) saveAsWebP(img image.Image, outputPath string, quality float32) error {
	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Configure WebP encoder options
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, quality)
	if err != nil {
		return fmt.Errorf("failed to create encoder options: %w", err)
	}

	// Encode to WebP
	if err := webp.Encode(outputFile, img, options); err != nil {
		return fmt.Errorf("failed to encode WebP: %w", err)
	}

	// Get file size for reporting
	fileInfo, err := outputFile.Stat()
	if err == nil {
		fmt.Printf("Output file size: %.2f KB\n", float64(fileInfo.Size())/1024)
	}

	return nil
}

// loadSVG loads and rasterizes an SVG file to an image
func (p *Processor) loadSVG(path string) (image.Image, error) {
	// Read SVG file
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open SVG file: %w", err)
	}
	defer file.Close()

	// Read file content
	svgData, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read SVG file: %w", err)
	}

	fmt.Printf("Loading SVG: %d bytes\n", len(svgData))

	// Parse SVG
	icon, err := oksvg.ReadIconStream(strings.NewReader(string(svgData)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse SVG: %w", err)
	}

	// Set reasonable dimensions with higher resolution for quality
	width, height := icon.ViewBox.W, icon.ViewBox.H
	if width == 0 || height == 0 {
		width, height = 512, 512 // Default size for SVGs without dimensions
	}

	// Scale up for better quality, then we'll scale down during cropping
	scale := 2.0
	renderWidth := int(width * scale)
	renderHeight := int(height * scale)

	fmt.Printf("Rendering SVG at %.0fx%.0f (2x resolution for quality)\n", width, height)

	// Create raster image with transparent background at higher resolution
	img := image.NewRGBA(image.Rect(0, 0, renderWidth, renderHeight))

	// Initialize with transparent background
	for y := 0; y < renderHeight; y++ {
		for x := 0; x < renderWidth; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 0}) // Transparent
		}
	}

	// Create scanner and rasterize at higher resolution
	scanner := rasterx.NewScannerGV(renderWidth, renderHeight, img, img.Bounds())
	raster := rasterx.NewDasher(renderWidth, renderHeight, scanner)

	// Set viewbox and draw with scaling
	icon.SetTarget(0, 0, width*scale, height*scale)
	icon.Draw(raster, 1.0)

	// Scale down for final image if we scaled up
	if scale != 1.0 {
		finalImg := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
		// Simple downsampling
		for y := 0; y < int(height); y++ {
			for x := 0; x < int(width); x++ {
				srcX := int(float64(x) * scale)
				srcY := int(float64(y) * scale)
				if srcX < renderWidth && srcY < renderHeight {
					finalImg.Set(x, y, img.At(srcX, srcY))
				}
			}
		}
		img = finalImg
	}

	fmt.Printf("SVG successfully converted to raster image\n")
	return img, nil
}

// loadSVGWithScale loads and rasterizes an SVG file with configurable scaling
func (p *Processor) loadSVGWithScale(path string, scale float32) (image.Image, error) {
	// Read SVG file
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open SVG file: %w", err)
	}
	defer file.Close()

	// Read file content
	svgData, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read SVG file: %w", err)
	}

	fmt.Printf("Loading SVG: %d bytes\n", len(svgData))

	// Parse SVG
	icon, err := oksvg.ReadIconStream(strings.NewReader(string(svgData)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse SVG: %w", err)
	}

	// Set reasonable dimensions with configurable resolution for quality
	width, height := icon.ViewBox.W, icon.ViewBox.H
	if width == 0 || height == 0 {
		width, height = 512, 512 // Default size for SVGs without dimensions
	}

	// Clamp scale between 1 and 4 for reasonable performance
	if scale < 1 {
		scale = 1
	}
	if scale > 4 {
		scale = 4
	}

	renderWidth := int(width * float64(scale))
	renderHeight := int(height * float64(scale))

	fmt.Printf("Rendering SVG at %.0fx%.0f (%.1fx scale for quality)\n", width, height, scale)

	// Create raster image with transparent background at higher resolution
	img := image.NewRGBA(image.Rect(0, 0, renderWidth, renderHeight))

	// Initialize with transparent background
	for y := 0; y < renderHeight; y++ {
		for x := 0; x < renderWidth; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 0}) // Transparent
		}
	}

	// Create scanner and rasterize at higher resolution
	scanner := rasterx.NewScannerGV(renderWidth, renderHeight, img, img.Bounds())
	raster := rasterx.NewDasher(renderWidth, renderHeight, scanner)

	// Set viewbox and draw with scaling
	icon.SetTarget(0, 0, width*float64(scale), height*float64(scale))
	icon.Draw(raster, 1.0)

	// Scale down for final image if we scaled up
	if scale != 1.0 {
		finalImg := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
		// Simple downsampling
		for y := 0; y < int(height); y++ {
			for x := 0; x < int(width); x++ {
				srcX := int(float64(x) * float64(scale))
				srcY := int(float64(y) * float64(scale))
				if srcX < renderWidth && srcY < renderHeight {
					finalImg.Set(x, y, img.At(srcX, srcY))
				}
			}
		}
		img = finalImg
	}

	fmt.Printf("SVG successfully converted to raster image\n")
	return img, nil
}

// LoadImage loads an image from file, supporting various formats including SVG
func (p *Processor) LoadImage(path string) (image.Image, error) {
	// Check if it's an SVG file
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".svg" {
		return p.loadSVG(path)
	}

	// Handle regular image formats
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Try to decode based on file extension
	switch ext {
	case ".png":
		return png.Decode(file)
	case ".jpg", ".jpeg":
		return jpeg.Decode(file)
	default:
		// Try generic decode
		img, _, err := image.Decode(file)
		return img, err
	}
}

// ScaleImage scales an image using various methods and saves it
func (p *Processor) ScaleImage(inputPath, outputPath string, factor float32, width, height int, algorithm string, quality, svgScale float32) error {
	var img image.Image
	var err error

	// Load the image
	if strings.ToLower(filepath.Ext(inputPath)) == ".svg" {
		img, err = p.loadSVGWithScale(inputPath, svgScale)
		if err != nil {
			return fmt.Errorf("failed to load SVG: %w", err)
		}
	} else {
		img, err = p.LoadImage(inputPath)
		if err != nil {
			return fmt.Errorf("failed to load image: %w", err)
		}
	}

	// Get original dimensions
	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	var targetWidth, targetHeight int

	// Calculate target dimensions
	if factor != 0 {
		// Scale by factor
		targetWidth = int(float32(originalWidth) * factor)
		targetHeight = int(float32(originalHeight) * factor)
	} else if width != 0 && height != 0 {
		// Both dimensions specified
		targetWidth = width
		targetHeight = height
	} else if width != 0 {
		// Only width specified, maintain aspect ratio
		aspectRatio := float32(originalHeight) / float32(originalWidth)
		targetWidth = width
		targetHeight = int(float32(width) * aspectRatio)
	} else if height != 0 {
		// Only height specified, maintain aspect ratio
		aspectRatio := float32(originalWidth) / float32(originalHeight)
		targetHeight = height
		targetWidth = int(float32(height) * aspectRatio)
	}

	fmt.Printf("Scaling from %dx%d to %dx%d\n", originalWidth, originalHeight, targetWidth, targetHeight)

	// Scale the image using the specified algorithm
	var scaledImg image.Image
	switch strings.ToLower(algorithm) {
	case "nearest":
		scaledImg = imaging.Resize(img, targetWidth, targetHeight, imaging.NearestNeighbor)
	case "bilinear", "linear":
		scaledImg = imaging.Resize(img, targetWidth, targetHeight, imaging.Linear)
	case "bicubic", "cubic":
		scaledImg = imaging.Resize(img, targetWidth, targetHeight, imaging.CatmullRom)
	case "lanczos":
		scaledImg = imaging.Resize(img, targetWidth, targetHeight, imaging.Lanczos)
	default:
		return fmt.Errorf("unsupported resampling algorithm: %s (use: nearest, bilinear, bicubic, lanczos)", algorithm)
	}

	// Save the scaled image
	return p.saveImage(scaledImg, outputPath, quality)
}

// saveImage saves an image in the appropriate format based on file extension
func (p *Processor) saveImage(img image.Image, outputPath string, quality float32) error {
	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Save based on file extension
	ext := strings.ToLower(filepath.Ext(outputPath))
	switch ext {
	case ".png":
		err = png.Encode(outputFile, img)
	case ".jpg", ".jpeg":
		// Convert quality from 0-100 to JPEG quality
		jpegQuality := int(quality)
		if jpegQuality < 1 {
			jpegQuality = 90
		}
		if jpegQuality > 100 {
			jpegQuality = 100
		}
		err = jpeg.Encode(outputFile, img, &jpeg.Options{Quality: jpegQuality})
	case ".webp":
		// Use WebP encoder
		options, encErr := encoder.NewLossyEncoderOptions(encoder.PresetDefault, quality)
		if encErr != nil {
			return fmt.Errorf("failed to create WebP encoder options: %w", encErr)
		}
		err = webp.Encode(outputFile, img, options)
	default:
		return fmt.Errorf("unsupported output format: %s (use: .png, .jpg, .jpeg, .webp)", ext)
	}

	if err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	// Get file size for reporting
	fileInfo, err := outputFile.Stat()
	if err == nil {
		fmt.Printf("Output file size: %.2f KB\n", float64(fileInfo.Size())/1024)
	}

	return nil
}
