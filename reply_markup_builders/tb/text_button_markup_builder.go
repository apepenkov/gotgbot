package tb

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// M - markup
func M(rows ...[]tgbotapi.KeyboardButton) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.ReplyKeyboardMarkup{
		Keyboard:       rows,
		ResizeKeyboard: true,
	}
}

// R - row
func R(buttons ...tgbotapi.KeyboardButton) []tgbotapi.KeyboardButton {
	return buttons
}

// B - button
func B(text, data string) tgbotapi.KeyboardButton {
	return tgbotapi.NewKeyboardButton(text)
}
