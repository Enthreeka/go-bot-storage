package callback

import (
	"github.com/Enthreeka/go-bot-storage/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BotSendTextCell(log *logger.Logger, bot *tgbotapi.BotAPI, userID int64) {
	text := "Напишите, какой раздел вы хотите добавить:"
	msg := tgbotapi.NewMessage(userID, text)

	_, err := bot.Send(msg)
	if err != nil {
		log.Error("failed to send message in CallbackQuery [create_cell] %v", err)
	}
}
