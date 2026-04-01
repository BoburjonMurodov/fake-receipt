package generator

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGenerateReceipt(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filename))

	gen := NewGenerator(
		filepath.Join(projectRoot, "templates"),
		filepath.Join(projectRoot, "fonts"),
	)

	templates := GetTemplates()
	tmpl, ok := templates["uzum_bank"]
	if !ok {
		t.Fatal("uzum_bank template not found")
	}

	data := map[string]string{
		"amount":        "55 000",
		"receiver_name": "Odilov Sarvar",
		"card_last4":    "4829",
	}

	imgBytes, err := gen.GenerateReceipt(tmpl, data)
	if err != nil {
		t.Fatalf("GenerateReceipt failed: %v", err)
	}

	if len(imgBytes) == 0 {
		t.Fatal("Generated image is empty")
	}

	// Save to /tmp for visual inspection
	outPath := "/tmp/test_receipt.png"
	err = os.WriteFile(outPath, imgBytes, 0644)
	if err != nil {
		t.Fatalf("Failed to write test image: %v", err)
	}

	t.Logf("Receipt saved to %s (%d bytes) — open to visually inspect", outPath, len(imgBytes))
}
