# Fake Receipt Telegram Bot - AI Context

Welcome, fellow agent! 🤖 This file contains everything you need to know about this repository so you can jump right in without having to read every file.

## Project Overview
This is a **Go**-based Telegram bot that generates realistic "fake" banking receipts. 
Users interact with the bot to select a bank template (e.g., Uzum Bank), provide details (amount, name, card number), and the bot returns a customized screenshot using the `fogleman/gg` rendering library.

## Tech Stack
- **Language:** Go
- **Telegram Library:** `github.com/go-telegram-bot-api/telegram-bot-api/v5`
- **Image Rendering:** `github.com/fogleman/gg` (for 2D text rendering & compositing)

## Architecture
The repository is split into a few focused packages:
- `main.go` — The entrypoint. Loads config and starts the bot listener.
- `config/config.go` — Loads environment variables. Requires `BOT_TOKEN`.
- `bot/handler.go` — Manages Telegram updates and the conversational state. The bot uses inline keyboards for template selection and stores user step progression in memory (`sessions` map).
- `generator/generator.go` — Core image manipulation logic. It loads a base PNG, paints over ("blanks out") original variable text areas using precise RGB values, and overlays new text using provided `.ttf` fonts. Contains `getAutoValue` logic for dynamic fields (like current iOS time).
- `generator/templates.go` — The single source of truth for template configurations. Defines the coordinates (X, Y), blanking rectangles, fonts, and colors for every variable text field on every template.

## How to Add a New Template
If the user asks you to add a new bank template (like Payme, Click, or Humo), follow these steps:

1. **Add Base Image**: Ask the user for a screenshot of the receipt and place it in the `templates/` directory as a PNG. Use the original app screenshot as the base.
2. **Define Template Options**: Go to `generator/templates.go` and create a new `*Template` function (e.g., `paymeTemplate()`). Add it to the `GetTemplates()` map.
3. **Map Coordinates**: 
   - Define exact X/Y parameters, colors, and `BlankRect` dimensions for each field (amount, name, card number, time).
   - Use the exact background color for `BlankColor` to seamlessly paint over original text.
   - For fields that should populate automatically (like the iOS status bar time), use `Auto: "time"`.
4. **Test & Adjust**: Before using it in the bot, run the unit test: `go test ./generator/... -v -run TestGenerateReceipt`. View the generated `/tmp/test_receipt.png` to tweak pixel coordinates, as they usually require 2-3 iterations of fine-tuning to look perfect.

## Bot Setup / Running
```bash
BOT_TOKEN=your_token go run main.go
```
