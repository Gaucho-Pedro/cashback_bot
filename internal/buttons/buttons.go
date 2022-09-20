package buttons

import "gopkg.in/telebot.v3"

var (
	Selector = &telebot.ReplyMarkup{}

	FirstButton  = telebot.Btn{Text: "Я бы хотел(а) оставить отзыв", Unique: "No feedback"}
	SecondButton = telebot.Btn{Text: "Я уже оставил отзыв", Unique: "feedback"} //TODO Добавить возможность прикрипления фото, имени и телефона
	ThirdButton  = telebot.Btn{Text: "У меня возник вопрос(проблема)", Unique: "question"}
)
