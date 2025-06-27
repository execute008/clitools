# CLI Tools

A collection of command-line utilities for various development tasks, with a focus on web optimization.

## Features

### Image Optimization

The image optimization tool helps prepare images for web use by:

1. **Cropping transparent areas** - Removes all transparent padding/offset to reduce file size
2. **WebP conversion** - Converts images to WebP format for better compression
3. **Quality control** - Adjustable quality settings for optimal size/quality balance

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd clitools
```

2. Install dependencies:
```bash
go mod download
```

3. Build the CLI tool:
```bash
go build -o clitools .
```

## Usage

### Image Optimization

Optimize an image by cropping transparent areas and converting to WebP:

```bash
# Basic usage
./clitools image optimize input.png output.webp

# With custom quality (0-100, default: 80)
./clitools image optimize input.jpg output.webp --quality 90

# Short flag version
./clitools image optimize input.png output.webp -q 75
```

**Supported input formats:**
- PNG
- JPEG/JPG
- Other formats supported by Go's image package

**Output format:**
- WebP (always, regardless of input format)

### Examples

```bash
# Optimize a PNG with transparency
./clitools image optimize logo.png logo-optimized.webp

# Optimize a JPEG photo
./clitools image optimize photo.jpg photo-optimized.webp --quality 85

# The tool will automatically add .webp extension if not provided
./clitools image optimize input.png output
# Output will be saved as: output.webp
```

## Web Optimization Benefits

This tool is specifically designed for web optimization:

1. **Smaller file sizes** - WebP format provides 25-35% better compression than JPEG/PNG
2. **Faster loading** - Removing transparent padding reduces unnecessary data
3. **Better performance** - Optimized images improve page load times
4. **Maintained quality** - Adjustable quality settings preserve visual fidelity

## Development

### Project Structure

```
clitools/
├── main.go                    # Entry point
├── cmd/                       # CLI commands
│   ├── root.go               # Root command
│   └── image.go              # Image manipulation commands
├── internal/
│   └── image/
│       └── processor.go      # Image processing logic
└── go.mod                    # Go module file
```

### Adding New Tools

To add a new tool to the CLI:

1. Create a new command file in `cmd/`
2. Add the command to `cmd/root.go`
3. Implement the logic in `internal/`

### Building for Distribution

```bash
# Build for current platform
go build -o clitools .

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o clitools-linux-amd64 .
GOOS=windows GOARCH=amd64 go build -o clitools-windows-amd64.exe .
GOOS=darwin GOARCH=amd64 go build -o clitools-darwin-amd64 .
```

## Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [go-webp](https://github.com/kolesa-team/go-webp) - WebP encoding/decoding

## License

[Add your license here]
