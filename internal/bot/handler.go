package bot

import (
	storage "bot/internal/storage"
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UpdateHandle struct {
	Bot     *tgapi.BotAPI
	Storage *storage.Storage
	Ans     Answers
}

func NewHandler(bot *tgapi.BotAPI, Storage *storage.Storage) *UpdateHandle {
	return &UpdateHandle{
		Bot:     bot,
		Storage: Storage,
		Ans:     GetAnswer(),
	}
}

func (m *UpdateHandle) Handle(upd tgapi.Update) {
	if upd.Message == nil {
		return
	}
	if upd.Message.IsCommand() {
		// Обработка команд
		switch upd.Message.Command() {
		case "run":
			m.CommandRun(upd.Message)
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
			if user.Action == "waiting" {
				m.Send(upd.Message.Chat.ID, m.Ans.AllReadyWaiting)
			} else if user.Action == "chat" {
				m.Send(user.PartnerID, upd.Message.Text)
			}
		}

	}
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
		m.Send(message.Chat.ID, m.Ans.BotNotRunned)
	} else {
		//Если пользователя есть в списке
		if user.Action == "waiting" {
			m.Send(message.Chat.ID, m.Ans.AllReadyWaiting)
		} else if user.Action == "chat" {
			m.Send(message.Chat.ID, m.Ans.Next)
			m.Send(user.PartnerID, m.Ans.ExitPartner)
			m.Storage.CleanPartner(message.Chat.ID)
			m.Storage.CleanPartner(user.PartnerID)
			m.StartChat(message.Chat.ID)
			m.StartChat(user.PartnerID)
		}
	}

}

func (m *UpdateHandle) CommandStart(message *tgapi.Message) {
	//Отправка сообщения о начале работы
	m.SendWithKeyboard(message.Chat.ID, m.Ans.StartFirst, []string{"/run"})
}

func (m *UpdateHandle) CommandRun(message *tgapi.Message) {
	//Добавление пользователя в список
	err := m.Storage.SetUser(message.Chat.ID, message.From.UserName)
	if err != nil {
		//Если пользователя уже в списке
		log.Println(err)
		m.Send(message.Chat.ID, m.Ans.AllReadyChatting)
		return
	}

	//Отправка сообщения о начале поиска
	m.SendWithCleanKeyboard(message.Chat.ID, m.Ans.StartSearch)

	//Начало чата
	m.StartChat(message.Chat.ID)
}
