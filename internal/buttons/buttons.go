package buttons

import "gopkg.in/telebot.v3"

var (
	selector = &telebot.ReplyMarkup{}

	firstButton  = selector.Data("Я бы хотел(а) оставить отзыв", "No_feedback")
	secondButton = selector.Data("Я уже оставил отзыв", "feedback") //TODO Добавить возможность прикрипления фото, имени и телефона
	thirdButton  = selector.Data("У меня возник вопрос(проблема)", "question")
)

type name struct {
}
