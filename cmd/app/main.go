package main

import (
	"cashback_bot/internal/config"
	"fmt"
	nestedFormatter "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages/memory"
	"gopkg.in/telebot.v3"
	"regexp"
	"time"
)

var (
	QuestionSG        = fsm.NewStateGroup("question")
	WaitQuestionState = QuestionSG.New("print")

	InputSG               = fsm.NewStateGroup("Input")
	InputPhotoState       = InputSG.New("photo")
	InputNameState        = InputSG.New("name")
	InputPartNumberState  = InputSG.New("partNumber")
	InputPhoneNumberState = InputSG.New("phoneNumber")

	firstButton    = telebot.Btn{Text: "Я бы хотел(а) оставить отзыв", Unique: "NoFeedback"}
	secondButton   = telebot.Btn{Text: "Я уже оставил отзыв", Unique: "feedback"}
	thirdButton    = telebot.Btn{Text: "У меня возник вопрос(проблема)", Unique: "question"}
	mainMenuButton = telebot.Btn{Text: "Главное меню"}
)

func main() {
	log.SetFormatter(&nestedFormatter.Formatter{
		NoColors:        true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	level, err := log.ParseLevel(config.GetConfig().LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)

	bot, err := telebot.NewBot(telebot.Settings{
		Token:     config.GetConfig().BotToken,
		Poller:    &telebot.LongPoller{Timeout: 20 * time.Second}, // TODO: Подумать над новым поллером с фильтром
		ParseMode: telebot.ModeMarkdown,
		Verbose:   config.GetConfig().BotDebug,
	})
	if err != nil {
		log.Fatal(err)
	}

	storage := memory.NewStorage()
	defer storage.Close()

	manager := fsm.NewManager(bot.Group(), storage)

	manager.Bind("/start", fsm.AnyState, OnStart(firstButton, secondButton, thirdButton, mainMenuButton))
	manager.Bind(&mainMenuButton, fsm.AnyState, OnMainMenu(firstButton, secondButton, thirdButton))
	bot.Handle(&firstButton, OnWantToFeedBack(secondButton))

	manager.Bind(&secondButton, fsm.DefaultState, OnFeedBackExist)
	manager.Bind(&thirdButton, fsm.DefaultState, OnQuestion)

	manager.Bind(telebot.OnText, WaitQuestionState, OnPrintQuestion())
	manager.Bind(telebot.OnText, fsm.DefaultState, OnAnswerFromAdmin(bot)) //TODO подумать над состоянием

	manager.Bind(telebot.OnPhoto, InputPhotoState, OnInputPhoto)
	manager.Bind(telebot.OnText, InputPhotoState, OnInputPhoto)
	manager.Bind(telebot.OnDocument, InputPhotoState, OnInputPhoto)

	manager.Bind(telebot.OnText, InputNameState, OnInputName)
	manager.Bind(telebot.OnText, InputPartNumberState, OnInputPartNumber)
	manager.Bind(telebot.OnText, InputPhoneNumberState, OnInputPhoneNumber(bot))

	bot.Start()
}

func OnStart(firstButton, secondButton, thirdButton, mainMenuButton telebot.Btn) fsm.Handler {
	menu := &telebot.ReplyMarkup{}
	menu.Inline(
		menu.Row(firstButton),
		menu.Row(secondButton),
		menu.Row(thirdButton),
	)
	menu.ResizeKeyboard = true

	mainMenu := &telebot.ReplyMarkup{}
	mainMenu.Reply(
		mainMenu.Row(mainMenuButton))
	mainMenu.ResizeKeyboard = true

	return func(context telebot.Context, state fsm.FSMContext) error {
		log.Infof("New user with id: %d", context.Chat().ID)
		state.Finish(true)
		context.Send("Спасибо, что выбрали SHIMA!\n\n"+
			"Хотим сделать Вам кешбек в размере 100 руб на телефон или карту.\n\n"+
			"Для получения кешбека Вам необходимо оставить отзыв о нашем продукте на Wildberries.\n", mainMenu)
		return context.Send("Подскажите, может Вы уже успели оставить отзыв?", menu)
	}
}

func OnMainMenu(firstButton, secondButton, thirdButton telebot.Btn) fsm.Handler {
	menu := &telebot.ReplyMarkup{}
	menu.Inline(
		menu.Row(firstButton),
		menu.Row(secondButton),
		menu.Row(thirdButton),
	)
	menu.ResizeKeyboard = true

	return func(context telebot.Context, state fsm.FSMContext) error {
		log.Debugf("[%d]: main menu", context.Chat().ID)
		state.Finish(true)
		return context.Send("Главное меню", menu)
	}
}

func OnWantToFeedBack(button telebot.Btn) telebot.HandlerFunc {
	menu := &telebot.ReplyMarkup{}
	menu.Inline(
		menu.Row(button),
	)
	return func(context telebot.Context) error {
		log.Debugf("[%d]: %s", context.Chat().ID, "wants to write a review")
		return context.Send(
			"Отлично!\n"+
				"Это очень простые 7 шагов :\n"+
				"1) Зайдите в Личный кабинет.\n"+
				"2) Зайдите в раздел \"Покупки\"\n"+
				"3) Выберите товар, который вы приобрели у нас\n"+
				"4) Кликните на \"Отзыв\" - > \"Оставить отзыв\"\n"+
				"5) Напишите, чем Вам понравился наш товар, прикрепите его фотографию, поставте оценку\n"+
				"6) Кликните \"Опубликовать отзыв\"\n"+
				"7) Сделайте скриншот готового отзыва и прикрепите в наш чат-бот.\n", menu)
	}
}

func OnQuestion(context telebot.Context, state fsm.FSMContext) error {
	log.Debugf("[%d]: %s", context.Chat().ID, "has a question")
	state.Set(WaitQuestionState)
	return context.Send("Опишите проблему с которой Вы столкнулись. Мы ответим на все Ваши вопросы и решим проблемы👌")
}

func OnPrintQuestion() fsm.Handler {
	return func(context telebot.Context, state fsm.FSMContext) error {
		err := context.ForwardTo(telebot.ChatID(config.GetConfig().AdminChatID))
		if err != nil {
			return err
		}
		//state.Set(fsm.DefaultState)
		state.Finish(true)
		return context.Send("Ваше сообщение принято")
	}
}

func OnAnswerFromAdmin(bot *telebot.Bot) fsm.Handler {
	return func(context telebot.Context, state fsm.FSMContext) error {
		if context.Message().IsReply() && context.Chat().ID == config.GetConfig().AdminChatID {
			log.Debug(context.Message().ReplyTo.OriginalSender.ID)
			_, err := bot.Copy(context.Message().ReplyTo.OriginalSender, context.Message())
			return err
		} else {
			return nil
		}
	}
}

func OnFeedBackExist(context telebot.Context, state fsm.FSMContext) error {
	log.Debugf("[%d]: %s", context.Chat().ID, "The review exists")
	state.Set(InputPhotoState)
	return context.Send("Класс!\n" +
		"Теперь нам необходим скриншот отзыва, ваше имя, артикул товара и номер телефона\n\n" +
		"Отправте, пожалуйста, скриншот вашего отзыва сюда")
}
func OnInputPhoto(context telebot.Context, state fsm.FSMContext) error {
	if context.Update().Message.Photo != nil {
		log.Debug(context.Update().Message.Photo.File.FileSize)
		state.Update("messageWithPhoto", context.Message())
		state.Set(InputNameState)
		return context.Send("Введите имя, под которым вы оставили отзыв на Wildberries\n" +
			"(Только **Имя**, без фамилии)")
	} else if context.Update().Message.Document != nil {
		return context.Send("Пожалуйста, пришлите мне фото как \"Фото\", а не как \"Файл\".")
	} else {
		return context.Send("Отправте, пожалуйста, фото")
	}
}

func OnInputName(context telebot.Context, state fsm.FSMContext) error {
	state.Update("name", context.Message().Text)
	state.Set(InputPartNumberState)
	return context.Send("Введите артикул товара (9 цифр)\n\n" +
		"Посмотреть его можно в личном кабинете WB.\n" +
		"Зайдите в раздел \"Профиль\" - >\"Покупки\"\n" +
		"Нажмите на товар, чуть ниже вы найдете артикул")
}

func OnInputPartNumber(context telebot.Context, state fsm.FSMContext) error {
	matched, _ := regexp.MatchString("^[0-9]{9}$", context.Message().Text)
	if matched {
		state.Update("article", "https://www.wildberries.ru/catalog/"+context.Message().Text+"/detail.aspx")
		state.Set(InputPhoneNumberState)
		return context.Send("Введите номер телефона")
	} else {
		return context.Send("Укажите верный артикул (9 цифр)")
	}
}

func OnInputPhoneNumber(bot *telebot.Bot) fsm.Handler {
	return func(context telebot.Context, state fsm.FSMContext) error {
		matched, _ := regexp.MatchString("^((\\+7|7|8)+([0-9]){10})$", context.Message().Text)
		if matched {
			_, err := bot.Send(telebot.ChatID(config.GetConfig().AdminChatID), fmt.Sprintf("Новый Отзыв!\n\nИмя: %s\nТелефон: %s\nСссылка на товар:%s\nСкриншот отзыва👇", state.MustGet("name"), context.Message().Text, state.MustGet("article")))
			if err != nil {
				log.Error(err)
			}
			_, err = bot.Forward(telebot.ChatID(config.GetConfig().AdminChatID), state.MustGet("messageWithPhoto").(*telebot.Message))
			if err != nil {
				log.Error(err)
			}
			state.Finish(true)
			return context.Send("Отлично! Наш менеджер проверит ваш отзыв и отправит кешбэк в течении 6 часов")
		} else {
			return context.Send("Укажите верный номер телефона")
		}
	}
}
