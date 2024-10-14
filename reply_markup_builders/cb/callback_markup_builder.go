package cb

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// this is made just for ease of creating inline keyboard markup and to make those tgbotapi.NewInlineKeyboardMarkup, tgbotapi.NewInlineKeyboardButtonData calls smaller

// M - markup
func M(rows ...[]tgbotapi.InlineKeyboardButton) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// R - row
func R(buttons ...tgbotapi.InlineKeyboardButton) []tgbotapi.InlineKeyboardButton {
	return buttons
}

// B - button
func B(text, data string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(text, data)
}

func Single(text, data string) tgbotapi.InlineKeyboardMarkup {
	return M(R(B(text, data)))
}

func Double(text1, data1, text2, data2 string) tgbotapi.InlineKeyboardMarkup {
	return M(R(B(text1, data1), B(text2, data2)))
}

var NONE = tgbotapi.InlineKeyboardMarkup{
	InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
}

func example() {
	// example of usage

	_ = M(
		R(
			B("button 1", "data 1"),
			B("button 2", "data 2"),
		),
		R(
			B("button 3", "data 3"),
			B("button 4", "data 4"),
		),
	)
	// the same as
	_ = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("button 1", "data 1"),
			tgbotapi.NewInlineKeyboardButtonData("button 2", "data 2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("button 3", "data 3"),
			tgbotapi.NewInlineKeyboardButtonData("button 4", "data 4"),
		),
	)
}
