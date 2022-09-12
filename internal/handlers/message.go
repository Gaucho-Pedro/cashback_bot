package handlers

import (
	"cashback_bot/internal/buttons"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func MessageHandler(message *tgBotApi.Message, bot *tgBotApi.BotAPI, cache map[int64]bool, adminChatID int64) {

	log.Debugf("[%s] %s", message.From.UserName, message.Text)
	// ChatID мой 461181622

	if message.Text != "" && cache[message.Chat.ID] {
		//msg := tgBotApi.NewForward(5481815893, message.Chat.ID, message.MessageID)

		bot.Send(tgBotApi.NewForward(adminChatID, message.Chat.ID, message.MessageID))
		msg := tgBotApi.NewMessage(message.Chat.ID, "Ваше сообщение принято")
		bot.Send(msg)
		cache[message.Chat.ID] = false
	}
	switch message.Command() {
	case "start":
		msg := tgBotApi.NewMessage(message.Chat.ID, "Спасибо, что выбрали SHIMA!\n\nХотим сделать Вам кешбек в размере 100 руб на телефон или вашу карту.\n\nДля получения кешбека Вам будет необходимо оставить отзыв о нашем продукте на WB.\nПодскажите, может Вы уже успели оставить отзыв?")
		msg.ParseMode = "markdown"
		msg.ReplyMarkup = buttons.MainKeyboard()
		bot.Send(msg)
	}
	log.Debugf("%s %v", "кеш:", cache)
	return
}

func ReplyHandler(message *tgBotApi.Message, bot *tgBotApi.BotAPI) {
	log.Debug(message.ReplyToMessage.ForwardFrom)
	if message.ReplyToMessage.ForwardFrom != nil {
		chatId := message.ReplyToMessage.ForwardFrom.ID
		log.Debug(chatId)
		bot.Send(tgBotApi.NewCopyMessage(chatId, message.Chat.ID, message.MessageID))
		return
	}
}
