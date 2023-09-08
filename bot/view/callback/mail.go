package callback

import (
	"github.com/Enthreeka/go-bot-storage/bot/controller"
	"github.com/Enthreeka/go-bot-storage/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Удалить 	cellController controller.Cell
type callbackMail struct {
	cellController controller.Cell

	bot *tgbotapi.BotAPI
	log *logger.Logger
}

func NewCallbackMail(cellController controller.Cell, bot *tgbotapi.BotAPI, log *logger.Logger) *callbackMail {
	return &callbackMail{
		cellController: cellController,
		bot:            bot,
		log:            log,
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
