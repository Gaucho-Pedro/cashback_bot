package main

import (
	"cashback_bot/internal/config"
	"cashback_bot/internal/handlers"
	nestedFormatter "github.com/antonfisher/nested-logrus-formatter"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&nestedFormatter.Formatter{
		NoColors:        true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	config := config.GetConfig()

	level, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)

	bot, err := tgBotApi.NewBotAPI(config.BotToken)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = config.BotDebug

	log.Infof("Authorized on account %s", bot.Self.UserName)

	u := tgBotApi.NewUpdate(0)
	u.Timeout = 60

	//updates := bot.GetUpdatesChan(u)
	cache := map[int64]bool{}
	for update := range bot.GetUpdatesChan(u) {
		if update.Message != nil && update.Message.ReplyToMessage == nil {
			go handlers.MessageHandler(update.Message, bot, cache, config.AdminChatID)
		} else if update.CallbackQuery != nil {
			go handlers.CallbackHandler(update.CallbackQuery, bot, cache)
		} else if update.Message.ReplyToMessage != nil && update.FromChat().ID == config.AdminChatID {
			go handlers.ReplyHandler(update.Message, bot)
		}
	}
}
