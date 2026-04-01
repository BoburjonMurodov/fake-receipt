package generator

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"time"

	"github.com/fogleman/gg"
)

// Generator handles receipt image generation
type Generator struct {
	templatesDir string
	fontsDir     string
}

// NewGenerator creates a new Generator instance
func NewGenerator(templatesDir, fontsDir string) *Generator {
	return &Generator{
		templatesDir: templatesDir,
		fontsDir:     fontsDir,
	}
}

// GenerateReceipt creates a receipt image by overlaying text on the template
func (g *Generator) GenerateReceipt(tmpl *Template, data map[string]string) ([]byte, error) {
	// Load base template image
	imgPath := filepath.Join(g.templatesDir, tmpl.BaseImage)
	baseImg, err := loadPNG(imgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load template image: %w", err)
	}

	bounds := baseImg.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	// Create drawing context from the base image
	dc := gg.NewContext(w, h)
	dc.DrawImage(baseImg, 0, 0)

	// Process each field
	for fieldName, fieldDef := range tmpl.Fields {
		value, ok := data[fieldName]
		if !ok {
			// Check if this is an auto field
			if fieldDef.Auto != "" {
				value = getAutoValue(fieldDef.Auto)
			} else {
				continue
			}
		}

		// Format the display value
		displayValue := formatFieldValue(fieldName, value)

		// Paint over the original text area with background color
		if fieldDef.BlankRect != [4]float64{} {
			r, gr, b, a := fieldDef.BlankColor.RGBA()
			dc.SetRGBA(
				float64(r)/65535.0,
				float64(gr)/65535.0,
				float64(b)/65535.0,
				float64(a)/65535.0,
			)
			dc.DrawRectangle(
				fieldDef.BlankRect[0],
				fieldDef.BlankRect[1],
				fieldDef.BlankRect[2],
				fieldDef.BlankRect[3],
			)
			dc.Fill()
		}

		// Load font
		fontPath := filepath.Join(g.fontsDir, fieldDef.FontFile)
		if err := dc.LoadFontFace(fontPath, fieldDef.FontSize); err != nil {
			return nil, fmt.Errorf("failed to load font %s: %w", fieldDef.FontFile, err)
		}

		// Set text color
		r, gr, b, a := fieldDef.Color.RGBA()
		dc.SetRGBA(
			float64(r)/65535.0,
			float64(gr)/65535.0,
			float64(b)/65535.0,
			float64(a)/65535.0,
		)

		// Draw text with alignment
		ax := 0.0 // left align
		switch fieldDef.Align {
		case "center":
			ax = 0.5
		case "right":
			ax = 1.0
		}
		dc.DrawStringAnchored(displayValue, fieldDef.X, fieldDef.Y, ax, 0.5)
	}

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, dc.Image()); err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	return buf.Bytes(), nil
}

// getAutoValue returns an auto-generated value based on the auto type
func getAutoValue(autoType string) string {
	switch autoType {
	case "time":
		return time.Now().Format("15:04")
	default:
		return ""
	}
}

// formatFieldValue formats the value for display on the receipt
func formatFieldValue(fieldName, value string) string {
	switch fieldName {
	case "amount":
		return fmt.Sprintf("Transferred %s sum", value)
	case "card_last4":
		return fmt.Sprintf("• %s", value)
	default:
		return value
	}
}

// loadPNG loads a PNG image from the given path
func loadPNG(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}
