package generator

import (
	"image/color"
)

// TextField defines how a text field should be rendered on the template
type TextField struct {
	X        float64     // X position
	Y        float64     // Y position
	FontSize float64     // Font size in points
	FontFile string      // Font filename (e.g. "Inter-Bold.ttf")
	Color    color.Color // Text color
	Align    string      // "left", "center", "right"
	// BlankRect defines the area to paint over before drawing text
	// Format: [x, y, width, height]
	BlankRect  [4]float64
	BlankColor color.Color
	// Auto indicates this field is auto-generated (e.g. "time")
	// Auto fields are not prompted to the user
	Auto string
}

// Template defines a receipt template with its base image and text fields
type Template struct {
	Name      string               // Display name (e.g. "Uzum Bank")
	Slug      string               // Identifier (e.g. "uzum_bank")
	BaseImage string               // Filename of base template image
	Fields    map[string]TextField // Field definitions
	// FieldOrder defines the order in which fields are prompted
	FieldOrder []string
	// FieldPrompts defines the prompt text shown to the user for each field
	FieldPrompts map[string]string
}

// Colors
var (
	colorWhite = color.RGBA{255, 255, 255, 255}
	colorGreen = color.RGBA{120, 255, 80, 255}
	// Uzum Bank purple background
	colorUzumPurple = color.RGBA{102, 16, 223, 255}
)

// GetTemplates returns all available receipt templates
func GetTemplates() map[string]*Template {
	return map[string]*Template{
		"uzum_bank": uzumBankTemplate(),
	}
}

func uzumBankTemplate() *Template {
	return &Template{
		Name:      "Uzum Bank",
		Slug:      "uzum_bank",
		BaseImage: "uzum_bank.png",
		FieldOrder: []string{"amount", "receiver_name", "card_last4"},
		FieldPrompts: map[string]string{
			"amount":        "💰 Enter the transfer amount (e.g. 56 000):",
			"receiver_name": "👤 Enter the receiver's full name:",
			"card_last4":    "💳 Enter the last 4 digits of the card:",
		},
		Fields: map[string]TextField{
			"time": {
				X:          92,
				Y:          102,
				FontSize:   48,
				FontFile:   "Inter-Bold.ttf",
				Color:      colorWhite,
				Align:      "left",
				BlankRect:  [4]float64{60, 64, 260, 56},
				BlankColor: colorUzumPurple,
				Auto:       "time",
			},
			"amount": {
				X:          80,
				Y:          540,
				FontSize:   88,
				FontFile:   "Inter-Bold.ttf",
				Color:      colorWhite,
				Align:      "left",
				BlankRect:  [4]float64{0, 460, 1320, 120},
				BlankColor: colorUzumPurple,
			},
			"receiver_name": {
				X:          80,
				Y:          680,
				FontSize:   84,
				FontFile:   "Inter-Bold.ttf",
				Color:      colorWhite,
				Align:      "left",
				BlankRect:  [4]float64{0, 590, 1320, 120},
				BlankColor: colorUzumPurple,
			},
			"card_last4": {
				X:          80,
				Y:          800,
				FontSize:   76,
				FontFile:   "Inter-Bold.ttf",
				Color:      colorWhite,
				Align:      "left",
				BlankRect:  [4]float64{0, 720, 600, 110},
				BlankColor: colorUzumPurple,
			},
		},
	}
}
