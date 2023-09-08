package callback

import (
	"github.com/Enthreeka/go-bot-storage/bot/view"
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

func (c *callbackMail) BotSendMainMenu(update *tgbotapi.Update) {
	userID := update.CallbackQuery.Message.Chat.ID

	msg := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    userID,
			MessageID: update.CallbackQuery.Message.MessageID,
		}, Text: "Управление разделами",
	}

	msg.ReplyMarkup = &view.StartKeyboard

	_, err := c.bot.Send(msg)
	if err != nil {
		c.log.Error("error sending main menu from callback: %v", err)
	}
}
