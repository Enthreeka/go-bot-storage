package command

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

func (c *cellView) ShowCell(update *tgbotapi.Update, msg *tgbotapi.MessageConfig) error {
	var userID int64
	if update.Message != nil && update.Message.Chat != nil {
		userID = update.Message.Chat.ID
	} else if update.CallbackQuery != nil && update.CallbackQuery.Message != nil && update.CallbackQuery.Message.Chat != nil {
		userID = update.CallbackQuery.Message.Chat.ID
	}

	cells, err := c.cellController.GetCell(userID)
	if err != nil {
		c.log.Error("failed to get cell: %v", err)
		return err
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	if len(cells) > 0 {
		for _, el := range cells {
			button := tgbotapi.NewInlineKeyboardButtonData(el.Name, fmt.Sprintf("cell_%s_%d", el.Name, el.ID))
			row = append(row, button)
			rows = append(rows, row)
			row = []tgbotapi.InlineKeyboardButton{}

		}

		rows = append(rows, row)
		rows = append(rows, []tgbotapi.InlineKeyboardButton{view.CreateCellButtonData, view.DeleteCellButtonData})

	} else {
		rows = append(rows, []tgbotapi.InlineKeyboardButton{view.CreateCellButtonData, view.DeleteCellButtonData})
	}
	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg.ParseMode = tgbotapi.ModeHTML
	msg.Text = "<b>Управление разделами</b>"
	msg.ReplyMarkup = &markup

	_, err = c.bot.Send(msg)
	if err != nil {
		c.log.Error("failed to send message in /start %v", err)
	}

	return nil
}

func (c *cellView) CreateCell(update *tgbotapi.Update, msg *tgbotapi.MessageConfig) error {
	if !view.KeyboardValidation(update.Message.Text) {
		c.log.Error("invalid button data")

		msg.Text = "Недопустимое название кнопки!"
		_, err := c.bot.Send(msg)
		if err != nil {
			c.log.Error("failed to send message in CreateCell %v", err)
			return err
		}
		return err
	}

	err := c.cellController.CreateCell(update)
	if err != nil {
		c.log.Error("failed create new cell by [%s]: %v", update.Message.From.UserName, err)

		msg.Text = "Ошибка при создании новой ячейки!"
		_, err = c.bot.Send(msg)
		if err != nil {
			c.log.Error("failed to send message in CreateCell %v", err)
		}
		return err
	}

	msg.Text = "Ячейка добавлена успешно"
	_, err = c.bot.Send(msg)
	if err != nil {
		c.log.Error("failed to send message in CreateCell %v", err)
	}

	return nil
}

func (c *cellView) CreateUnderCell(update *tgbotapi.Update, msg *tgbotapi.MessageConfig, cellData *string) error {
	if !view.KeyboardValidation(update.Message.Text) {
		c.log.Error("invalid button data")

		msg.Text = "Недопустимое название кнопки!"
		_, err := c.bot.Send(msg)
		if err != nil {
			c.log.Error("failed to send message in CreateUnderCell %v", err)
			return err
		}
		return err
	}

	cellID, _ := model.FindIdName(*cellData)

	err := c.cellController.CreateUnderCell(update, &cellID)
	if err != nil {
		c.log.Error("failed create new [create_under_cell] by [%s]: %v", update.Message.From.UserName, err)

		msg.Text = "Ошибка при создании новой темы!"
		_, err = c.bot.Send(msg)
		if err != nil {
			c.log.Error("failed to send message in CreateUnderCell %v", err)
		}
		return err
	}

	msg.ParseMode = tgbotapi.ModeHTML
	msg.Text = fmt.Sprintf("Тема: <b>%s</b> добавлена успешно", update.Message.Text)

	_, err = c.bot.Send(msg)
	if err != nil {
		c.log.Error("failed to send message in CreateUnderCell %v", err)
	}
	return nil
}

func (c *cellView) ShowUnderCell(update *tgbotapi.Update, cellData *string) error {
	userID := update.Message.Chat.ID

	msg := tgbotapi.NewMessage(userID, "")
	msg.ParseMode = tgbotapi.ModeHTML

	cellID, name := model.FindIdName(*cellData)

	underCell, err := c.cellController.GetUnderCell(userID, cellID)
	if err != nil {
		c.log.Error("failed to show underCell [create_under_cell] by [%s]: %v", update.Message.From.UserName, err)
		return err
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

	msg.ReplyMarkup = &markup
	msg.Text = fmt.Sprintf("<b>%s</b>", name)

	_, err = c.bot.Send(msg)
	if err != nil {
		c.log.Error("failed to send message after add create_under_cell %v", err)
	}

	return nil
}
