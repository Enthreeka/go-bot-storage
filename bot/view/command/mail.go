package command

import (
	"github.com/Enthreeka/go-bot-storage/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type commandMail struct {
	bot *tgbotapi.BotAPI
	log *logger.Logger
}

func NewCommandMail(bot *tgbotapi.BotAPI, log *logger.Logger) *commandMail {
	return &commandMail{
		bot: bot,
		log: log,
	}
}

func (c *commandMail) BotSendInfo(msg *tgbotapi.MessageConfig) {
	msg.Text = "Данный бот реализует хранилище ваших ссылок, pdf файлов и заметок. \n" +
		"Вы можете создать ячейки с общими темами, внутри уже которых создавать ячейки для определенных тем."

	_, err := c.bot.Send(msg)
	if err != nil {
		c.log.Error("failed to send message in /info %v", err)
	}
}

func (c *commandMail) BotSendDefault(msg *tgbotapi.MessageConfig) {
	msg.Text = "Неверная команда, попробуйте /start"

	_, err := c.bot.Send(msg)
	if err != nil {
		c.log.Error("failed to send message in default %v", err)
	}
}
