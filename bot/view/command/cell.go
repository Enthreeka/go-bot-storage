package command

import (
	"github.com/Enthreeka/go-bot-storage/bot/controller"
	"github.com/Enthreeka/go-bot-storage/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type cellView struct {
	cellController controller.Cell

	bot *tgbotapi.BotAPI
	log *logger.Logger
}

func NewCellView(cellController controller.Cell, bot *tgbotapi.BotAPI, log *logger.Logger) *cellView {
	return &cellView{
		cellController: cellController,
		bot:            bot,
		log:            log,
	}
}

func (c *cellView) CreateCell(update *tgbotapi.Update, msg *tgbotapi.MessageConfig) error {
	err := c.cellController.CreateCell(update)
	if err != nil {
		c.log.Error("failed create new cell by [%s]: %v", update.Message.From.UserName, err)

		msg.Text = "Ошибка при создании новой ячейки!"
		c.bot.Send(msg)
		if err != nil {
			c.log.Error("failed to send message in CreateCell %v", err)
		}
	}

	msg.Text = "Ячейка добавлена успешно"
	c.bot.Send(msg)
	if err != nil {
		c.log.Error("failed to send message in CreateCell %v", err)
	}

	return nil
}
