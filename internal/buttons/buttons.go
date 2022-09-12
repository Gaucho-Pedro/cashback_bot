package buttons

import tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func MainKeyboard() *tgBotApi.InlineKeyboardMarkup {
	keyboard := tgBotApi.NewInlineKeyboardMarkup(
		tgBotApi.NewInlineKeyboardRow(tgBotApi.NewInlineKeyboardButtonData("Я бы хотел(а) оставить отзыв", "No feedback")),
		tgBotApi.NewInlineKeyboardRow(tgBotApi.NewInlineKeyboardButtonData("Я уже оставил отзыв", "feedback")),
		tgBotApi.NewInlineKeyboardRow(tgBotApi.NewInlineKeyboardButtonData("У меня возник вопрос(проблема)", "question")))
	return &keyboard
}
