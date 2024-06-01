package bot

type Answers struct {
	StartFirst       string
	StartSearch      string
	StartChat        string
	Next             string
	Exit             string
	ExitPartner      string
	Help             string
	AllReadyChatting string
	AllReadyWaiting  string
	BotNotRunned     string
	ZeroUser         string
}

func GetAnswer() Answers {
	return Answers{
		StartFirst:       "👋 Привет! Я бот, который помогает найти собеседника. Напиши /run, чтобы начать!",
		StartSearch:      "🔍 Ищу для тебя собеседника, подожди немного...",
		StartChat:        "🎉 Ура! Нашел тебе собеседника, начнем общение!",
		Next:             "🔄 Беседа окончена. Ищу для тебя нового собеседника...",
		Exit:             "👋 Ты вышел из чата. Если захочешь вернуться, всегда рад!",
		ExitPartner:      "🚪 Собеседник ушел. Ищу для тебя нового друга...",
		Help:             "ℹ️ Чтобы начать общение, напиши /run и следуй инструкциям!",
		AllReadyChatting: "💬 Ты уже в чате! Если хочешь нового собеседника, напиши /next.",
		AllReadyWaiting:  "⏳ Ты уже в очереди на поиск. Жди немного!",
		BotNotRunned:     "🚀 Ты еще не начал общение. Напиши /run, чтобы начать!",
		ZeroUser:         "🕐 Пока собеседников нет, но ты в очереди. Немного подожди!",
	}

}
