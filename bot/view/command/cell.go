package command

import (
	"github.com/Enthreeka/go-bot-storage/bot/controller"
	"github.com/Enthreeka/go-bot-storage/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type cellView struct {
	cellController controller.Cell

	log *logger.Logger
}

func NewCellView(cellController controller.Cell, log *logger.Logger) *cellView {
	return &cellView{
		cellController: cellController,
		log:            log,
	}
}

func (c *cellView) CreateCell(update *tgbotapi.Update, bot *tgbotapi.BotAPI, msg *tgbotapi.MessageConfig) error {
	err := c.cellController.CreateCell(update)
	if err != nil {
		c.log.Error("failed create new cell by [%s]: %v", update.Message.From.UserName, err)

		msg.Text = "Ошибка при создании новой ячейки!"
		bot.Send(msg)
		if err != nil {
			c.log.Error("failed to send message in CreateCell %v", err)
		}
	}

	msg.Text = "Ячейка добавлена успешно"
	bot.Send(msg)
	if err != nil {
		c.log.Error("failed to send message in CreateCell %v", err)
	}

	return nil
}
