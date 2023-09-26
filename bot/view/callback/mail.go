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

func (c *callbackMail) BotSendTextDeleteCell(userID int64) {
	text := "Нажмите на раздел, который вы хотите УДАЛИТЬ.\nВ случае, если вы передумали, вызовите команду /start"
	msg := tgbotapi.NewMessage(userID, text)

	_, err := c.bot.Send(msg)
	if err != nil {
		c.log.Error("failed to send message in CallbackQuery [delete_cell] %v", err)
	}
}

func (c *callbackMail) BotSendTextDeleteUnderCell(userID int64) {
	text := "Нажмите на тему, которую вы хотите УДАЛИТЬ.\nВ случае, если вы передумали, вызовите команду /start"
	msg := tgbotapi.NewMessage(userID, text)

	_, err := c.bot.Send(msg)
	if err != nil {
		c.log.Error("failed to send message in CallbackQuery [delete_under_cell] %v", err)
	}
}

func (c *callbackMail) BotSendTextUpdateData(userID int64) {
	text := "Отправьте новые данные, которые будут расположены в данной теме.\nВ случае, если вы передумали, вызовите команду /start"
	msg := tgbotapi.NewMessage(userID, text)

	_, err := c.bot.Send(msg)
	if err != nil {
		c.log.Error("failed to send message in CallbackQuery [update_data] %v", err)
	}
}
