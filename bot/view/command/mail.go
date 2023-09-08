package command

import (
	"github.com/Enthreeka/go-bot-storage/bot/controller"
	"github.com/Enthreeka/go-bot-storage/bot/view"
	"github.com/Enthreeka/go-bot-storage/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type commandMail struct {
	cellController controller.Cell

	bot *tgbotapi.BotAPI
	log *logger.Logger
}

func NewCommandMail(cellController controller.Cell, bot *tgbotapi.BotAPI, log *logger.Logger) *commandMail {
	return &commandMail{
		cellController: cellController,
		bot:            bot,
		log:            log,
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

func (c *commandMail) BotSendStart(msg *tgbotapi.MessageConfig) {
	msg.Text = "Управление разделами"
	msg.ReplyMarkup = view.StartKeyboard

	_, err := c.bot.Send(msg)
	if err != nil {
		c.log.Error("failed to send message in /start %v", err)
	}
}
