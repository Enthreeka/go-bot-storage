package callback

import (
	"fmt"
	"github.com/Enthreeka/go-bot-storage/bot/controller"
	"github.com/Enthreeka/go-bot-storage/bot/model"
	"github.com/Enthreeka/go-bot-storage/bot/view"
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
		button := tgbotapi.NewInlineKeyboardButtonData(el.Name, fmt.Sprintf("cell_%s_%d", el.Name, el.ID))
		row = append(row, button)
		rows = append(rows, row)
		row = []tgbotapi.InlineKeyboardButton{}

	}

	rows = append(rows, row)
	rows = append(rows, []tgbotapi.InlineKeyboardButton{view.MainMenuButtonData})

	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    userID,
			MessageID: update.CallbackQuery.Message.MessageID,
		}, Text: "Главный раздел",
	}

	msg.ReplyMarkup = &markup

	_, err = c.bot.Send(msg)
	if err != nil {
		c.log.Error("error sending cell keyboard: %v", err)
	}

	return nil
}

func (c *cellView) ShowUnderCell(update *tgbotapi.Update) (int, error) {
	userID := update.CallbackQuery.Message.Chat.ID

	_, name := model.FindIntStr(update.CallbackQuery.Data)

	underCell, err := c.cellController.GetUnderCell(userID)
	if err != nil {
		return 0, err
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton
	if len(underCell) > 0 {
		for _, el := range underCell {
			button := tgbotapi.NewInlineKeyboardButtonData(el.Name, fmt.Sprintf("underCell_%s_%d", el.Name, el.ID))
			row = append(row, button)
			rows = append(rows, row)
			row = []tgbotapi.InlineKeyboardButton{}

		}
		rows = append(rows, row)

		rows = append(rows, []tgbotapi.InlineKeyboardButton{view.CellButtonDataCreate})
		rows = append(rows, []tgbotapi.InlineKeyboardButton{view.CellButtonDataDelete})
		rows = append(rows, []tgbotapi.InlineKeyboardButton{view.MainMenuButtonData})

	} else {
		rows = append(rows, []tgbotapi.InlineKeyboardButton{view.NoOneRowsButtonData})
		rows = append(rows, []tgbotapi.InlineKeyboardButton{view.CellButtonDataCreate})
		rows = append(rows, []tgbotapi.InlineKeyboardButton{view.MainMenuButtonData})
	}

	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    userID,
			MessageID: update.CallbackQuery.Message.MessageID,
		}, Text: name,
	}

	msg.ReplyMarkup = &markup

	_, err = c.bot.Send(msg)
	if err != nil {
		c.log.Error("error sending under cell keyboard: %v", err)
	}

	return 0, nil

}

func (c *cellView) CreateUnderCell(update *tgbotapi.Update, cellID *int) error {
	err := c.cellController.CreateUnderCell(update, cellID)
	if err != nil {
		return err
	}
	return nil
}
