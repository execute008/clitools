package cmd

import (
	"clitools/internal/image"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Image manipulation tools for web optimization",
	Long: `Image manipulation tools that help optimize images for web use.
Includes cropping transparent areas and converting to WebP format.`,
}

var optimizeCmd = &cobra.Command{
	Use:   "optimize [input] [output]",
	Short: "Optimize image by cropping transparent areas and converting to WebP",
	Long: `Optimize an image for web use by:
1. Cropping all transparent padding/offset areas
2. Converting to WebP format for better compression
3. Reducing file size while maintaining quality

Examples:
  clitools image optimize input.png output.webp
  clitools image optimize image.jpg optimized.webp`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]
		outputPath := args[1]

		// Ensure output has .webp extension
		if !strings.HasSuffix(strings.ToLower(outputPath), ".webp") {
			ext := filepath.Ext(outputPath)
			if ext != "" {
				outputPath = strings.TrimSuffix(outputPath, ext) + ".webp"
			} else {
				outputPath = outputPath + ".webp"
			}
			fmt.Printf("Output will be saved as: %s\n", outputPath)
		}

		quality, _ := cmd.Flags().GetFloat32("quality")
		svgScale, _ := cmd.Flags().GetFloat32("svg-scale")

		processor := image.NewProcessor()
		err := processor.OptimizeImageWithScale(inputPath, outputPath, quality, svgScale)
		if err != nil {
			return fmt.Errorf("failed to optimize image: %w", err)
		}

		fmt.Printf("Successfully optimized %s -> %s\n", inputPath, outputPath)
		return nil
	},
}

var scaleCmd = &cobra.Command{
	Use:   "scale [input] [output]",
	Short: "Scale/resize images with various options",
	Long: `Scale or resize images with flexible options:
- Scale by factor: --factor 0.5 (50% smaller) or --factor 2.0 (2x larger)
- Set specific dimensions: --width 800 --height 600
- Fit to width/height maintaining aspect ratio: --width 800 or --height 600
- Resize with different resampling algorithms for quality

Examples:
  clitools image scale input.png output.png --factor 0.5
  clitools image scale input.jpg output.jpg --width 800 --height 600
  clitools image scale input.webp output.webp --width 1200
  clitools image scale input.svg output.png --height 400 --algorithm lanczos`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]
		outputPath := args[1]

		// Get flags
		factor, _ := cmd.Flags().GetFloat32("factor")
		width, _ := cmd.Flags().GetInt("width")
		height, _ := cmd.Flags().GetInt("height")
		algorithm, _ := cmd.Flags().GetString("algorithm")
		quality, _ := cmd.Flags().GetFloat32("quality")
		svgScale, _ := cmd.Flags().GetFloat32("svg-scale")

		// Validate parameters
		if factor == 0 && width == 0 && height == 0 {
			return fmt.Errorf("must specify either --factor, --width, or --height")
		}

		if factor != 0 && (width != 0 || height != 0) {
			return fmt.Errorf("cannot use --factor with --width or --height")
		}

		processor := image.NewProcessor()
		err := processor.ScaleImage(inputPath, outputPath, factor, width, height, algorithm, quality, svgScale)
		if err != nil {
			return fmt.Errorf("failed to scale image: %w", err)
		}

		fmt.Printf("Successfully scaled %s -> %s\n", inputPath, outputPath)
		return nil
	},
}

func init() {
	imageCmd.AddCommand(optimizeCmd)
	imageCmd.AddCommand(scaleCmd)

	// Add flags for optimize command
	optimizeCmd.Flags().Float32P("quality", "q", 80, "WebP quality (0-100)")
	optimizeCmd.Flags().Float32P("svg-scale", "s", 2, "SVG rendering scale factor for quality (1-4)")

	// Add flags for scale command
	scaleCmd.Flags().Float32P("factor", "f", 0, "Scale factor (e.g., 0.5 for 50%, 2.0 for 200%)")
	scaleCmd.Flags().IntP("width", "w", 0, "Target width in pixels")
	scaleCmd.Flags().Int("height", 0, "Target height in pixels")
	scaleCmd.Flags().StringP("algorithm", "a", "lanczos", "Resampling algorithm: nearest, bilinear, bicubic, lanczos")
	scaleCmd.Flags().Float32P("quality", "q", 90, "Output quality for JPEG/WebP (0-100)")
	scaleCmd.Flags().Float32P("svg-scale", "s", 2, "SVG rendering scale factor for quality (1-4)")
}
