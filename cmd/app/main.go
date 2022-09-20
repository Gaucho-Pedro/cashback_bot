package main

import (
	"cashback_bot/internal/config"
	nestedFormatter "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages/memory"
	"gopkg.in/telebot.v3"
	"time"
)

var (
	QuestionSG        = fsm.NewStateGroup("question")
	WaitQuestionState = QuestionSG.New("print")

	firstButton  = telebot.Btn{Text: "Я бы хотел(а) оставить отзыв", Unique: "NoFeedback"}
	secondButton = telebot.Btn{Text: "Я уже оставил отзыв", Unique: "feedback"} //TODO Добавить возможность прикрипления фото, имени и телефона
	thirdButton  = telebot.Btn{Text: "У меня возник вопрос(проблема)", Unique: "question"}
	//TODO Добавить проверку имени, артикула, номера телефона
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

	storage := memory.NewStorage()
	defer storage.Close()

	manager := fsm.NewManager(bot.Group(), storage)

	bot.Handle("/start", OnStart(firstButton, secondButton, thirdButton))

	bot.Handle(&firstButton, OnWantToFeedBack(secondButton))
	bot.Handle(&secondButton, OnFeedBackExist())
	manager.Bind(&thirdButton, fsm.DefaultState, OnQuestion)

	manager.Bind(telebot.OnText, WaitQuestionState, OnPrintQuestion(bot, config.AdminChatID))
	manager.Bind(telebot.OnText, fsm.AnyState, OnAnswerFromAdmin(bot, config.AdminChatID))

	bot.Start()
}

func OnStart(firstButton, secondButton, thirdButton telebot.Btn) telebot.HandlerFunc {
	menu := &telebot.ReplyMarkup{}
	menu.Inline(
		menu.Row(firstButton),
		menu.Row(secondButton),
		menu.Row(thirdButton),
	)
	menu.ResizeKeyboard = true
	return func(context telebot.Context) error {
		log.Debugf("New user with id: %d", context.Chat().ID)
		return context.Send("Спасибо, что выбрали SHIMA!\n\n"+
			"Хотим сделать Вам кешбек в размере 100 руб на телефон или вашу карту.\n\n"+
			"Для получения кешбека Вам будет необходимо оставить отзыв о нашем продукте на WB.\n"+
			"Подскажите, может Вы уже успели оставить отзыв?", menu)
	}
}

func OnWantToFeedBack(button telebot.Btn) telebot.HandlerFunc {
	menu := &telebot.ReplyMarkup{}
	menu.Inline(
		menu.Row(button),
	)
	return func(context telebot.Context) error {
		log.Debugf("[%d]: %s", context.Chat().ID, "wants to write a review")
		return context.Send("1) Зайдите в свой кабинет на WB\n"+
			"2) Кликните на \"Отзывы и вопросы\"\n"+
			"3) Сделайте скрин отзыва и пришлите в чат\n"+
			"4) В течение 6 часов мы вам пришлем кешбек.\n🙂", menu)
	}
}

func OnFeedBackExist() telebot.HandlerFunc {
	return func(context telebot.Context) error {
		log.Debugf("[%d]: %s", context.Chat().ID, "The review exists")
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
	}
}

func OnQuestion(context telebot.Context, state fsm.FSMContext) error {
	log.Debugf("[%d]: %s", context.Chat().ID, "has a question")
	state.Set(WaitQuestionState)
	return context.Send("Опишите проблему с которой Вы столкнулись. Мы ответим на все Ваши вопросы и решим проблемы👌")
}

func OnPrintQuestion(bot *telebot.Bot, adminChatID int64) fsm.Handler {
	return func(context telebot.Context, state fsm.FSMContext) error {
		_, err := bot.Forward(telebot.ChatID(adminChatID), context.Message())
		if err != nil {
			log.Error(err)
		}
		state.Set(fsm.DefaultState)
		return context.Send("Ваше сообщение принято")
	}
}

func OnAnswerFromAdmin(bot *telebot.Bot, adminChatID int64) fsm.Handler {
	return func(context telebot.Context, state fsm.FSMContext) error {
		if context.Message().IsReply() && context.Chat().ID == adminChatID {
			log.Debug(context.Message().ReplyTo.OriginalSender.ID)
			_, err := bot.Copy(context.Message().ReplyTo.OriginalSender, context.Message())
			return err
		} else {
			return nil
		}
	}
}
