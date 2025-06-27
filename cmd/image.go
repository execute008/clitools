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
		
		processor := image.NewProcessor()
		err := processor.OptimizeImage(inputPath, outputPath, quality)
		if err != nil {
			return fmt.Errorf("failed to optimize image: %w", err)
		}

		fmt.Printf("Successfully optimized %s -> %s\n", inputPath, outputPath)
		return nil
	},
}

func init() {
	imageCmd.AddCommand(optimizeCmd)
	
	// Add flags
	optimizeCmd.Flags().Float32P("quality", "q", 80, "WebP quality (0-100)")
}
