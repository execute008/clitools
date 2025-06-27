package image

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
)

type Processor struct{}

func NewProcessor() *Processor {
	return &Processor{}
}

// OptimizeImage crops transparent areas and converts to WebP format
func (p *Processor) OptimizeImage(inputPath, outputPath string, quality float32) error {
	// Open and decode input image
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
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

// LoadImage loads an image from file, supporting various formats
func (p *Processor) LoadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	// Try to decode based on file extension
	ext := strings.ToLower(filepath.Ext(path))
	
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
