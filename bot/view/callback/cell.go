package callback

import (
	"fmt"
	"github.com/Enthreeka/go-bot-storage/bot/controller"
	"github.com/Enthreeka/go-bot-storage/bot/model"
	"github.com/Enthreeka/go-bot-storage/bot/view"
	"github.com/Enthreeka/go-bot-storage/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
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
	rows = append(rows, []tgbotapi.InlineKeyboardButton{view.CreateCellButtonData, view.DeleteCellButtonData})

	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    userID,
			MessageID: update.CallbackQuery.Message.MessageID,
		},
		Text: "Главный раздел",
	}

	msg.ReplyMarkup = &markup

	_, err = c.bot.Send(msg)
	if err != nil {
		c.log.Error("error sending cell keyboard: %v", err)
	}

	return nil
}

func (c *cellView) ShowUnderCell(update *tgbotapi.Update, data string) (int, error) {
	userID := update.CallbackQuery.Message.Chat.ID

	var cellID int
	var name string
	if strings.HasPrefix(update.CallbackQuery.Data, "underCell_") {
		cellID, name = model.FindIdName(data)
	} else {
		cellID, name = model.FindIdName(update.CallbackQuery.Data)
	}

	underCell, err := c.cellController.GetUnderCell(userID, cellID)
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

		rows = append(rows, []tgbotapi.InlineKeyboardButton{view.CellButtonDataCreate, view.CellButtonDataDelete})
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
		c.log.Info("sending under cell keyboard -> create new message: %v", err)

		msg := tgbotapi.NewMessage(userID, "")
		msg.ReplyMarkup = &markup
		msg.Text = name

		_, err := c.bot.Send(msg)
		if err != nil {
			c.log.Error("failed to send message after delete under_cell %v", err)
		}

	}

	return cellID, nil
}

func (c *cellView) DeleteCell(update *tgbotapi.Update) (*tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")

	if !strings.HasPrefix(update.CallbackQuery.Data, "cell_") {
		msg.Text = "Для удаления был выбран не раздел!"
		_, err := c.bot.Send(msg)
		if err != nil {
			c.log.Error("failed to send message in DeleteCell %v", err)
		}

		return &msg, err
	}

	cellID, name := model.FindIdName(update.CallbackQuery.Data)
	err := c.cellController.DeleteCell(cellID)
	if err != nil {
		return &msg, err
	}

	msg.Text = fmt.Sprintf("Раздел: %s,удален успешно", name)
	_, err = c.bot.Send(msg)
	if err != nil {
		c.log.Error("failed to send message in DeleteCell %v", err)
	}

	if resp, err := c.bot.Request(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID)); nil != err || !resp.Ok {
		c.log.Error("failed to delete message id %d (%s): %v", update.CallbackQuery.Message.MessageID, string(resp.Result), err)
	}

	return &msg, err
}

func (c *cellView) DeleteUnderCell(update *tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")

	if !strings.HasPrefix(update.CallbackQuery.Data, "underCell_") {
		msg.Text = "Для удаления была выбрана не тема!"
		_, err := c.bot.Send(msg)
		if err != nil {
			c.log.Error("failed to send message in DeleteUnderCell %v", err)
		}

		return err
	}

	cellID, name := model.FindIdName(update.CallbackQuery.Data)
	err := c.cellController.DeleteUnderCell(cellID)
	if err != nil {
		return err
	}

	msg.Text = fmt.Sprintf("Раздел: %s,удален успешно", name)

	_, err = c.bot.Send(msg)
	if err != nil {
		c.log.Error("failed to send message in DeleteUnderCell %v", err)
	}

	//TODO вынести Request
	if resp, err := c.bot.Request(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID)); nil != err || !resp.Ok {
		c.log.Error("failed to delete message id %d (%s): %v", update.CallbackQuery.Message.MessageID, string(resp.Result), err)
	}

	return err
}
