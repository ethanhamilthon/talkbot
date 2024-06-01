package main

import (
	bot "bot/internal/bot"
)

func main() {
	b := bot.New()
	b.Start()
}
