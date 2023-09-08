package callback

import (
	"fmt"
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

func (c *cellView) ShowCell(update *tgbotapi.Update) error {
	userID := update.CallbackQuery.Message.Chat.ID

	cells, err := c.cellController.GetCell(userID)
	if err != nil {
		return err
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	for _, el := range cells {
		button := tgbotapi.NewInlineKeyboardButtonData(el.Name, fmt.Sprintf("%s_%d", el.Name, el.ID))
		row = append(row, button)
		rows = append(rows, row)
		row = []tgbotapi.InlineKeyboardButton{}

	}

	rows = append(rows, row)
	//rows = append(rows, []tgbotapi.InlineKeyboardButton{goToMainMenuKeyboard})

	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewEditMessageReplyMarkup(userID, update.CallbackQuery.Message.MessageID, markup)

	_, err = c.bot.Send(msg)
	if err != nil {
		c.log.Error("error sending cell keyboard: %v", err)
	}

	return nil
}
