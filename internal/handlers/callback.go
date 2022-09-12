package handlers

import (
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func CallbackHandler(callbackQuery *tgBotApi.CallbackQuery, bot *tgBotApi.BotAPI, cache map[int64]bool) {
	log.Debugf("[%s] %s", callbackQuery.From.UserName, callbackQuery.Data)
	log.Debug(callbackQuery.Message.Chat.ID)
	msg := tgBotApi.NewMessage(callbackQuery.Message.Chat.ID, "")

	msg.ParseMode = "markdown"

	switch callbackQuery.Data {
	case "No feedback":
		msg.Text = "Шаг 1. Зайдите в свой кабинет на WB\nШаг 2. Кликните на \"Отзывы и вопросы\"\nШаг 3. Сделайте скрин отзыва и пришлите в чат\nШаг 4. В течение 6 часов мы вам пришлем кешбек.\n🙂"
		bot.Send(msg)
	case "feedback":
		msg.Text = "Отлично! 😊\n\nНа данном этапе, для получения бонуса, Вам нужно оставить отзыв. Это очень простые 8 шагов :\n1️⃣Зайдите в Личный кабинет.\n2️⃣Найдите раздел “Покупки”\n3️⃣Выберите товар, который вы приобрели\n4️⃣Кликните на “Отзыв” - > “Оставить отзыв”\n5️⃣Напишите, чем Вам понравился наш товар, прикрепите его фотографию\n6️⃣Кликните “Опубликовать отзыв”\n7️⃣Сделайте скриншот готового отзыва и прикрепите в наш чат-бот.\n8.В течение 6 часов мы вам пришлем кэшбек."
		bot.Send(msg)
	case "question":
		msg.Text = "Опишите проблему с которой Вы столкнулись. Мы ответим на все Ваши вопросы и решим проблем👌"
		cache[callbackQuery.Message.Chat.ID] = true
		bot.Send(msg)
	}
	log.Debugf("%s %v", "кеш:", cache)
	return
}
