package bot

import (
	"fmt"
	"log"
	"sync"

	"fake-receipt/generator"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// UserSession tracks the state of a user's conversation
type UserSession struct {
	Step         int
	TemplateName string
	Data         map[string]string
}

// Handler manages the Telegram bot interactions
type Handler struct {
	bot       *tgbotapi.BotAPI
	gen       *generator.Generator
	templates map[string]*generator.Template
	sessions  map[int64]*UserSession
	mu        sync.RWMutex
}

// NewHandler creates a new bot handler
func NewHandler(bot *tgbotapi.BotAPI, gen *generator.Generator) *Handler {
	return &Handler{
		bot:       bot,
		gen:       gen,
		templates: generator.GetTemplates(),
		sessions:  make(map[int64]*UserSession),
	}
}

// Start begins listening for updates
func (h *Handler) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.bot.GetUpdatesChan(u)

	log.Printf("Bot started as @%s", h.bot.Self.UserName)

	for update := range updates {
		if update.CallbackQuery != nil {
			h.handleCallback(update.CallbackQuery)
			continue
		}

		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			h.handleCommand(update.Message)
		} else {
			h.handleTextInput(update.Message)
		}
	}
}

// handleCommand handles /start and other commands
func (h *Handler) handleCommand(msg *tgbotapi.Message) {
	switch msg.Command() {
	case "start":
		h.sendTemplateSelection(msg.Chat.ID)
	default:
		reply := tgbotapi.NewMessage(msg.Chat.ID, "Unknown command. Use /start to begin.")
		h.bot.Send(reply)
	}
}

// sendTemplateSelection shows available templates as inline buttons
func (h *Handler) sendTemplateSelection(chatID int64) {
	// Clear any existing session
	h.mu.Lock()
	delete(h.sessions, chatID)
	h.mu.Unlock()

	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, tmpl := range h.templates {
		btn := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("🧾 %s", tmpl.Name),
			fmt.Sprintf("template:%s", tmpl.Slug),
		)
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(btn))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)

	reply := tgbotapi.NewMessage(chatID,
		"🧾 *Fake Receipt Generator*\n\nSelect a receipt template:",
	)
	reply.ParseMode = "Markdown"
	reply.ReplyMarkup = keyboard
	h.bot.Send(reply)
}

// handleCallback handles inline keyboard button presses
func (h *Handler) handleCallback(cb *tgbotapi.CallbackQuery) {
	chatID := cb.Message.Chat.ID

	// Answer the callback to remove loading state
	callback := tgbotapi.NewCallback(cb.ID, "")
	h.bot.Request(callback)

	// Parse callback data
	if len(cb.Data) > 9 && cb.Data[:9] == "template:" {
		slug := cb.Data[9:]
		tmpl, ok := h.templates[slug]
		if !ok {
			reply := tgbotapi.NewMessage(chatID, "❌ Template not found.")
			h.bot.Send(reply)
			return
		}

		// Initialize session
		h.mu.Lock()
		h.sessions[chatID] = &UserSession{
			Step:         0,
			TemplateName: slug,
			Data:         make(map[string]string),
		}
		h.mu.Unlock()

		// Send first field prompt
		reply := tgbotapi.NewMessage(chatID,
			fmt.Sprintf("✅ Selected: *%s*\n\n%s", tmpl.Name, tmpl.FieldPrompts[tmpl.FieldOrder[0]]),
		)
		reply.ParseMode = "Markdown"
		h.bot.Send(reply)
	}
}

// handleTextInput processes user text input during the conversation flow
func (h *Handler) handleTextInput(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID

	h.mu.RLock()
	session, ok := h.sessions[chatID]
	h.mu.RUnlock()

	if !ok {
		reply := tgbotapi.NewMessage(chatID, "Use /start to begin creating a receipt.")
		h.bot.Send(reply)
		return
	}

	tmpl, ok := h.templates[session.TemplateName]
	if !ok {
		reply := tgbotapi.NewMessage(chatID, "❌ Session error. Use /start to begin again.")
		h.bot.Send(reply)
		return
	}

	// Store the current field value
	currentField := tmpl.FieldOrder[session.Step]
	session.Data[currentField] = msg.Text

	// Move to next step
	session.Step++

	// Check if we have all fields
	if session.Step >= len(tmpl.FieldOrder) {
		// Generate the receipt
		h.generateAndSend(chatID, tmpl, session.Data)

		// Clean up session
		h.mu.Lock()
		delete(h.sessions, chatID)
		h.mu.Unlock()
		return
	}

	// Ask for the next field
	nextField := tmpl.FieldOrder[session.Step]
	prompt := tmpl.FieldPrompts[nextField]
	reply := tgbotapi.NewMessage(chatID, prompt)
	h.bot.Send(reply)
}

// generateAndSend creates the receipt image and sends it to the user
func (h *Handler) generateAndSend(chatID int64, tmpl *generator.Template, data map[string]string) {
	// Send a "generating" message
	waitMsg := tgbotapi.NewMessage(chatID, "⏳ Generating your receipt...")
	h.bot.Send(waitMsg)

	// Generate the receipt image
	imgBytes, err := h.gen.GenerateReceipt(tmpl, data)
	if err != nil {
		log.Printf("Error generating receipt: %v", err)
		reply := tgbotapi.NewMessage(chatID, "❌ Failed to generate receipt. Please try again with /start")
		h.bot.Send(reply)
		return
	}

	// Send the image
	photoFile := tgbotapi.FileBytes{
		Name:  "receipt.png",
		Bytes: imgBytes,
	}
	photo := tgbotapi.NewPhoto(chatID, photoFile)
	photo.Caption = "🧾 Here's your receipt!"
	h.bot.Send(photo)

	// Offer to create another
	reply := tgbotapi.NewMessage(chatID, "Want to create another? Use /start")
	h.bot.Send(reply)
}
