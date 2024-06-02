package bot

import (
	cfg "bot/internal/config"
	storage "bot/internal/storage"
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	config *cfg.Config
}

func New() *Bot {
	return &Bot{
		config: cfg.New(),
	}
}

func (b *Bot) Start() {
	//Стартуем бота
	bot, err := tgapi.NewBotAPI(b.config.GetTelegramToken())
	if err != nil {
		panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	//Конфиги бота
	u := tgapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	//Создание хэндлера сообщений
	store := storage.New()
	updateHandle := NewHandler(bot, store)

	//Основной цикл работы бота
	for update := range updates {
		go updateHandle.Handle(update)
	}
}
