package bot

import (
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (m *MessageHandle) Send(ID int64, text string) {
	msg := tgapi.NewMessage(ID, text)
	m.Bot.Send(msg)
}

func (m *MessageHandle) SendWithKeyboard(ID int64, text string, keys []string) {
	msg := tgapi.NewMessage(ID, text)
	Buttons := make([]tgapi.KeyboardButton, 0)
	for _, key := range keys {
		Buttons = append(Buttons, tgapi.NewKeyboardButton(key))
	}
	msg.ReplyMarkup = tgapi.NewReplyKeyboard(Buttons)
	m.Bot.Send(msg)
}

func (m *MessageHandle) SendWithCleanKeyboard(ID int64, text string) {
	msg := tgapi.NewMessage(ID, text)
	msg.ReplyMarkup = tgapi.NewRemoveKeyboard(true)
	m.Bot.Send(msg)
}

func (m *MessageHandle) StartChat(ID int64) {
	//получение пользователя
	user, err := m.Storage.GetUser(ID)
	if err != nil {
		m.Send(ID, m.Ans.BotNotRunned)
		return
	}

	//Если пользователя уже в чате
	if user.Action == "chat" {
		m.Send(ID, m.Ans.AllReadyChatting)
		return
	}

	//Если пользователь в списке ожидания
	if user.Action == "waiting" {
		partner, err := m.Storage.GetPartner(ID)
		if err != nil && err.Error() == "No partner found" {
			m.Send(ID, m.Ans.ZeroUser)
			return
		}
		//Если есть собеседник
		m.Storage.CreateChat(ID, partner.ID)

		//Отправка сообщения о начале чата
		m.Send(ID, m.Ans.StartChat)
		m.Send(partner.ID, m.Ans.StartChat)
	}
}
