package main

import (
	"log"

	"fake-receipt/bot"
	"fake-receipt/config"
	"fake-receipt/generator"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize Telegram bot
	botAPI, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	log.Printf("Authorized on account %s", botAPI.Self.UserName)

	// Initialize image generator
	gen := generator.NewGenerator(cfg.TemplatesDir, cfg.FontsDir)

	// Initialize and start the bot handler
	handler := bot.NewHandler(botAPI, gen)
	handler.Start()
}
