package main

import (
	"cashback_bot/internal/config"
	nestedFormatter "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
	"time"
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

	bot, err := telebot.NewBot(telebot.Settings{
		Token:     config.BotToken,
		Poller:    &telebot.LongPoller{Timeout: 60 * time.Second},
		ParseMode: telebot.ModeMarkdown,
		//Verbose:   true,
	})
	if err != nil {
		log.Fatal(err)
	}
	cache := map[int64]bool{}

	selector1 := &telebot.ReplyMarkup{}
	selector2 := &telebot.ReplyMarkup{}

	firstButton := telebot.Btn{Text: "Я бы хотел(а) оставить отзыв", Unique: "No feedback"}
	secondButton := telebot.Btn{Text: "Я уже оставил отзыв", Unique: "feedback"} //TODO Добавить возможность прикрипления фото, имени и телефона
	thirdButton := telebot.Btn{Text: "У меня возник вопрос(проблема)", Unique: "question"}
	//TODO Добавить проверку имени, артикула, номера телефона

	selector1.Inline(
		selector1.Row(firstButton),
		selector1.Row(secondButton),
		selector1.Row(thirdButton),
	)

	selector2.Inline(
		selector2.Row(secondButton),
	)

	bot.Handle("/start", func(c telebot.Context) error {
		return c.Send("Спасибо, что выбрали SHIMA!\n\n"+
			"Хотим сделать Вам кешбек в размере 100 руб на телефон или вашу карту.\n\n"+
			"Для получения кешбека Вам будет необходимо оставить отзыв о нашем продукте на WB.\n"+
			"Подскажите, может Вы уже успели оставить отзыв?", selector1)
	})
	bot.Handle(&firstButton, func(context telebot.Context) error {
		return context.Send("1) Зайдите в свой кабинет на WB\n"+
			"2) Кликните на \"Отзывы и вопросы\"\n"+
			"3) Сделайте скрин отзыва и пришлите в чат\n"+
			"4) В течение 6 часов мы вам пришлем кешбек.\n🙂", selector2)
	})
	bot.Handle(&secondButton, func(context telebot.Context) error {
		log.Debug("2")
		return context.Send("Отлично! 😊\n\n" +
			"На данном этапе, для получения бонуса, Вам нужно оставить отзыв. Это очень простые 8 шагов :\n" +
			"1) Зайдите в Личный кабинет.\n" +
			"2) Найдите раздел “Покупки”\n" +
			"3) Выберите товар, который вы приобрели\n" +
			"4) Кликните на “Отзыв” - > “Оставить отзыв”\n" +
			"5) Напишите, чем Вам понравился наш товар, прикрепите его фотографию\n" +
			"6) Кликните “Опубликовать отзыв”\n" +
			"7) Сделайте скриншот готового отзыва и прикрепите в наш чат-бот.\n" +
			"8) В течение 6 часов мы вам пришлем кэшбек.")
	})
	bot.Handle(&thirdButton, func(context telebot.Context) error {
		log.Debug("3")
		cache[context.Chat().ID] = true
		return context.Send("Опишите проблему с которой Вы столкнулись. Мы ответим на все Ваши вопросы и решим проблемы👌")
	})
	bot.Handle(telebot.OnText, func(context telebot.Context) error {
		if !context.Message().IsReply() && cache[context.Chat().ID] {
			_, err := bot.Forward(telebot.ChatID(config.AdminChatID), context.Message())
			if err != nil {
				log.Error(err)
			}
			cache[context.Chat().ID] = false
			return context.Send("Ваше сообщение принято")
		} else if context.Message().IsReply() && context.Chat().ID == config.AdminChatID {
			log.Debug(context.Message().ReplyTo.OriginalSender.ID)
			_, err := bot.Copy(context.Message().ReplyTo.OriginalSender, context.Message())
			return err
		} else {
			return nil
		}
	})
	bot.Start()
}
