package bot

import "fmt"

type Answers struct {
	StartFirst,
	StartSearch,
	Next,
	Exit,
	ExitPartner,
	Help,
	AllReadyChatting,
	AllReadyWaiting,
	BotNotRunned,
	ZeroUser,
	Error,
	OfflinePartner string
	StartChat func(theme string) string
}

func GetAnswer() Answers {
	return Answers{
		StartFirst:  "👋 Привет! Я бот, который помогает найти собеседника. Напиши /run, чтобы начать!",
		StartSearch: "🔍 Ищу для тебя собеседника, подожди немного...",
		StartChat: func(theme string) string {
			return fmt.Sprintf("🎉 Ура! Нашел тебе собеседника, начнем общение!\nТема разговора:%s", theme)
		},
		Next:             "🔄 Беседа окончена. Ищу для тебя нового собеседника...",
		Exit:             "👋 Ты вышел из чата. Если захочешь вернуться, всегда рад!",
		ExitPartner:      "🚪 Собеседник ушел. Ищу для тебя нового друга...",
		Help:             "ℹ️ Чтобы начать общение, напиши /run и следуй инструкциям!",
		AllReadyChatting: "💬 Ты уже в чате! Если хочешь нового собеседника, напиши /next.",
		AllReadyWaiting:  "⏳ Ты уже в очереди на поиск. Жди немного!",
		BotNotRunned:     "🚀 Ты еще не начал общение. Напиши /run, чтобы начать!",
		ZeroUser:         "🕐 Пока собеседников нет, но ты в очереди. Немного подожди!",
		OfflinePartner:   "😴 Собеседник не в сети.",
		Error:            "🤕 Произошла ошибка. Попробуй еще раз!",
	}

}
