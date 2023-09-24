package callback

import (
	"github.com/Enthreeka/go-bot-storage/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type callbackMail struct {
	bot *tgbotapi.BotAPI
	log *logger.Logger
}

func NewCallbackMail(bot *tgbotapi.BotAPI, log *logger.Logger) *callbackMail {
	return &callbackMail{
		bot: bot,
		log: log,
	}
}

func (c *callbackMail) BotSendTextCell(userID int64) {
	text := "Напишите, какой раздел вы хотите добавить:"
	msg := tgbotapi.NewMessage(userID, text)

	_, err := c.bot.Send(msg)
	if err != nil {
		c.log.Error("failed to send message in CallbackQuery [create_cell] %v", err)
	}
}

func (c *callbackMail) BotSendTextUnderCell(userID int64) {
	text := "Напишите, какую темы вы хотите добавить:"
	msg := tgbotapi.NewMessage(userID, text)

	_, err := c.bot.Send(msg)
	if err != nil {
		c.log.Error("failed to send message in CallbackQuery [create_under_cell] %v", err)
	}
}

func (c *callbackMail) BotSendTextData(userID int64) {
	text := "Напишите, что вы хотите разместить в данной теме, либо прикрепите файл:"
	msg := tgbotapi.NewMessage(userID, text)

	_, err := c.bot.Send(msg)
	if err != nil {
		c.log.Error("failed to send message in CallbackQuery [add_data] %v", err)
	}
}
