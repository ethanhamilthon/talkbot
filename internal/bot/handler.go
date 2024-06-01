package bot

import (
	storage "bot/internal/storage"
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageHandle struct {
	Bot     *tgapi.BotAPI
	Storage *storage.Storage
	Ans     Answers
}

func NewHandler(bot *tgapi.BotAPI, Storage *storage.Storage) *MessageHandle {
	return &MessageHandle{
		Bot:     bot,
		Storage: Storage,
		Ans:     GetAnswer(),
	}
}

func (m *MessageHandle) Handle(message *tgapi.Message) {
	if message.IsCommand() {
		//Обработка команд
		if message.Command() == "run" {
			m.CommandRun(message)
		} else if message.Command() == "start" {
			m.CommandStart(message)
		} else if message.Command() == "next" {
			m.CommandNext(message)
		} else if message.Command() == "exit" {
			m.CommandExit(message)
		}
	} else {
		//Обработка сообщений
		user, err := m.Storage.GetUser(message.Chat.ID)
		if err != nil {
			//Если пользователя нет в списке
			m.Send(message.Chat.ID, m.Ans.BotNotRunned)
		} else {
			//Если пользователя есть в списке
			if user.Action == "waiting" {
				m.Send(message.Chat.ID, m.Ans.AllReadyWaiting)
			} else if user.Action == "chat" {
				m.Send(user.PartnerID, message.Text)
			}
		}

	}
}

func (m *MessageHandle) CommandExit(message *tgapi.Message) {
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

func (m *MessageHandle) CommandNext(message *tgapi.Message) {
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

func (m *MessageHandle) CommandStart(message *tgapi.Message) {
	//Отправка сообщения о начале работы
	m.SendWithKeyboard(message.Chat.ID, m.Ans.StartFirst, []string{"/run"})
}

func (m *MessageHandle) CommandRun(message *tgapi.Message) {
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
