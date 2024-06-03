package bot

import (
	"errors"
	"math/rand"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (m *UpdateHandle) Send(ID int64, text string) {
	msg := tgapi.NewMessage(ID, text)
	m.Bot.Send(msg)
}

type Keyboard struct {
	Type    string // "inline" or "menu"
	Buttons [][]string
}

func (m *UpdateHandle) SendWithKeyboard(ID int64, text string, keys Keyboard) {
	msg := tgapi.NewMessage(ID, text)
	switch keys.Type {
	case "inline":
		
	case "menu":
		Buttons := make([][]tgapi.KeyboardButton, 0, len(keys.Buttons))
		for _, row := range keys.Buttons {
			Row := make([]tgapi.KeyboardButton, 0, len(row))
			for _, btn := range row {
				Row = append(Row, tgapi.NewKeyboardButton(btn))
			}
			Buttons = append(Buttons, tgapi.NewKeyboardButtonRow(Row...))
		}
		msg.ReplyMarkup = tgapi.NewReplyKeyboard(Buttons...)
	}
	m.Bot.Send(msg)
}

func (m *UpdateHandle) SendWithCleanKeyboard(ID int64, text string) {
	msg := tgapi.NewMessage(ID, text)
	msg.ReplyMarkup = tgapi.NewRemoveKeyboard(true)
	m.Bot.Send(msg)
}

func (m *UpdateHandle) StartChat(ID int64) (int64, error) {
	//получение пользователя
	user, err := m.Storage.GetUser(ID)
	if err != nil {
		return 0, err
	}

	//Если пользователя уже в чате
	if user.Action == "chat" {
		return 0, errors.New("User is already in chat")
	}

	//Если пользователь в списке ожидания
	if user.Action == "waiting" {
		partner, err := m.Storage.GetPartner(user)
		if err != nil {
			return 0, err
		}
		//Если есть собеседник
		m.Storage.CreateChat(ID, partner.ID)

		return partner.ID, nil
	}

	return 0, errors.New("User is offline")
}

func (m *UpdateHandle) CloseChat(ID int64) (int64, error) {
	//получение пользователя
	user, err := m.Storage.GetUser(ID)
	if err != nil {
		return 0, err
	}

	//Удалить партнеров
	m.Storage.CleanPartner(ID)
	m.Storage.CleanPartner(user.PartnerID)

	return user.PartnerID, nil
}

func (m *UpdateHandle) GetRandomTheme() string {
	return m.Themes[rand.Intn(len(m.Themes))].ThemeText
}
