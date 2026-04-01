package config

import (
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	BotToken     string
	TemplatesDir string
	FontsDir     string
}

func Load() *Config {
	// Get project root relative to this file
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filename))

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		panic("BOT_TOKEN environment variable is required")
	}

	return &Config{
		BotToken:     botToken,
		TemplatesDir: filepath.Join(projectRoot, "templates"),
		FontsDir:     filepath.Join(projectRoot, "fonts"),
	}
}
