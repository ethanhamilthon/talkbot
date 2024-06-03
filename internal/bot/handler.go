package bot

import (
	storage "bot/internal/storage"
	"fmt"
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)



type UpdateHandle struct {
	Bot     *tgapi.BotAPI
	Storage *storage.Storage
	Ans     Answers
	Themes  []Theme
}

func NewHandler(bot *tgapi.BotAPI, Storage *storage.Storage) *UpdateHandle {
	return &UpdateHandle{
		Bot:     bot,
		Storage: Storage,
		Ans:     GetAnswer(),
		Themes:  GetThemes(),
	}
}

//Сообщений можно отправлять только в handle

func (m *UpdateHandle) Handle(upd tgapi.Update) {
	if upd.Message == nil {
		return
	}
	if upd.Message.IsCommand() {
		// Обработка команд
		switch upd.Message.Command() {
		case "run":
			m.CommandRun(upd.Message.Chat.ID, upd.Message.From.UserName)
		case "start":
			m.CommandStart(upd.Message)
		case "next":
			m.CommandNext(upd.Message)
		case "exit":
			m.CommandExit(upd.Message)
		}
	} else {
		//Обработка сообщений
		user, err := m.Storage.GetUser(upd.Message.Chat.ID)
		if err != nil {
			//Если пользователя нет в списке
			m.Send(upd.Message.Chat.ID, m.Ans.BotNotRunned)
		} else {
			//Если пользователя есть в списке
			switch user.Action {
			case "waiting":
				m.Send(upd.Message.Chat.ID, m.Ans.AllReadyWaiting)
			case "chat":
				m.SendText(user, upd.Message)
			case "offline":
				m.Send(upd.Message.Chat.ID, m.Ans.BotNotRunned)
			}
		}

	}
}

func (m *UpdateHandle) SendText(user storage.User, message *tgapi.Message) {
	ok := m.Storage.IsPartnerOnline(user.PartnerID)
	if !ok {
		m.Send(user.ID, m.Ans.OfflinePartner)
		m.Storage.CleanPartner(user.ID)
		m.StartChat(user.ID)
		return
	}
	text := fmt.Sprintf("💬%s: %s", user.PartnerName, message.Text)
	m.Send(user.PartnerID, text)
}

func (m *UpdateHandle) CommandExit(message *tgapi.Message) {
	//Получение пользователя
	user, err := m.Storage.GetUser(message.Chat.ID)
	if err != nil {
		//Если пользователя нет в списке
		m.Send(message.Chat.ID, m.Ans.BotNotRunned)
		return
	}
	//Удаление пользователя из списка
	m.Send(message.Chat.ID, m.Ans.Exit)
	m.Storage.DeleteUser(message.Chat.ID)

	//Поиск нового чата для собеседника
	m.Send(user.PartnerID, m.Ans.ExitPartner)
	m.Storage.CleanPartner(user.PartnerID)
	m.StartChat(user.PartnerID)
}

func (m *UpdateHandle) CommandNext(message *tgapi.Message) {
	//Получение пользователя
	user, err := m.Storage.GetUser(message.Chat.ID)
	if err != nil {
		//Если пользователя нет в списке
		m.Send(user.ID, m.Ans.BotNotRunned)
	} else {
		//Если пользователя есть в списке
		if user.Action == "waiting" {
			m.Send(user.ID, m.Ans.AllReadyWaiting)
		} else if user.Action == "chat" {
			m.Send(user.ID, m.Ans.Next)
			m.Send(user.ID, m.Ans.ExitPartner)
			m.CommandRun(user.ID, user.Name)
			m.CommandRun(user.PartnerID, user.PartnerName)
		}
	}

}

func (m *UpdateHandle) CommandStart(message *tgapi.Message) {
	//Отправка сообщения о начале работы
	keyboard := Keyboard{
		Type:    "menu",
		Buttons: [][]string{{"/run"}},
	}
	m.SendWithKeyboard(message.Chat.ID, m.Ans.StartFirst, keyboard)
}

func (m *UpdateHandle) CommandRun(ID int64, UserName string) {
	//Добавление пользователя в список
	err := m.Storage.SetUser(ID, UserName)
	if err != nil {
		//Если пользователя уже в списке
		log.Println(err)
		m.Send(ID, m.Ans.AllReadyChatting)
		return
	}

	//Отправка сообщения о начале поиска
	m.SendWithCleanKeyboard(ID, m.Ans.StartSearch)

	//Начало чата
	partnerID, err := m.StartChat(ID)
	if err != nil {
		switch err.Error() {
		case "User is already in chat":
			m.Send(ID, m.Ans.AllReadyChatting)
		case "User is offline":
			m.Send(ID, m.Ans.BotNotRunned)
		case "No partner found":
			m.Send(ID, m.Ans.ZeroUser)
		case "User does not exist":
			m.Send(ID, m.Ans.BotNotRunned)
		default:
			m.Send(ID, m.Ans.Error)
		}
		return
	}
	theme := m.GetRandomTheme()
	m.Send(ID, m.Ans.StartChat(theme))
	m.Send(partnerID, m.Ans.StartChat(theme))
}
