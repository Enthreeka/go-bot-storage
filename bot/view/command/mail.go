package command

import (
	"github.com/Enthreeka/go-bot-storage/bot/view"
	"github.com/Enthreeka/go-bot-storage/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BotSendInfo(log *logger.Logger, bot *tgbotapi.BotAPI, msg *tgbotapi.MessageConfig) {
	msg.Text = "Данный бот реализует хранилище ваших ссылок, pdf файлов и заметок. \n" +
		"Вы можете создать ячейки с общими темами, внутри уже которых создавать ячейки для определенных тем."

	_, err := bot.Send(msg)
	if err != nil {
		log.Error("failed to send message in /info %v", err)
	}
}

func BotSendDefault(log *logger.Logger, bot *tgbotapi.BotAPI, msg *tgbotapi.MessageConfig) {
	msg.Text = "Неверная команда, попробуйте /start"

	_, err := bot.Send(msg)
	if err != nil {
		log.Error("failed to send message in default %v", err)
	}
}

func BotSendStart(log *logger.Logger, bot *tgbotapi.BotAPI, msg *tgbotapi.MessageConfig) {
	msg.Text = "Управление разделами"
	msg.ReplyMarkup = view.StartKeyboard

	_, err := bot.Send(msg)
	if err != nil {
		log.Error("failed to send message in /start %v", err)
	}
}
