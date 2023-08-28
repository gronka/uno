package uf

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	BotId string
	//TODO: make a separate channel for different log levels
	ChannelId int64
}

// To get botm updates and ChatIds
//curl -s https://api.telegram.org/bot5911437983:AAFsKTavTMsLS2AJseQn5hrVwQnoJOfGzgI/getUpdates

var telegramBot TelegramBot

func (tger *TelegramBot) SendMsg(level, body string) {
	bot, err := tgbotapi.NewBotAPI(tger.BotId)
	if err != nil {
		LoggingError(err)
	}

	// Enabling this gives more info about the bot
	//bot.Debug = true

	msg := tgbotapi.NewMessage(tger.ChannelId, body)

	if _, err := bot.Send(msg); err != nil {
		LoggingError(err)
	}
}
