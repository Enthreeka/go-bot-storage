package command

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

func (c *cellView) ShowCell(update *tgbotapi.Update, msg *tgbotapi.MessageConfig) error {
	var userID int64
	if update.Message != nil && update.Message.Chat != nil {
		userID = update.Message.Chat.ID
	} else if update.CallbackQuery != nil && update.CallbackQuery.Message != nil && update.CallbackQuery.Message.Chat != nil {
		userID = update.CallbackQuery.Message.Chat.ID
	}

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

	//TODO проверка если в len(cell) == 0
	rows = append(rows, row)
	rows = append(rows, []tgbotapi.InlineKeyboardButton{view.CreateCellButtonData, view.DeleteCellButtonData})

	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg.Text = "Управление разделами"
	msg.ReplyMarkup = &markup

	_, err = c.bot.Send(msg)
	if err != nil {
		c.log.Error("failed to send message in /start %v", err)
	}

	return nil
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

func (c *cellView) CreateUnderCell(update *tgbotapi.Update, msg *tgbotapi.MessageConfig, cellID *int) error {
	err := c.cellController.CreateUnderCell(update, cellID)
	if err != nil {
		c.log.Error("failed create new under_cell by [%s]: %v", update.Message.From.UserName, err)

		msg.Text = "Ошибка при создании новой темы!"
		c.bot.Send(msg)
		if err != nil {
			c.log.Error("failed to send message in CreateUnderCell %v", err)
		}
	}

	msg.Text = fmt.Sprintf("Тема: %s ,добавлена успешно", update.Message.Text)
	_, err = c.bot.Send(msg)
	if err != nil {
		c.log.Error("failed to send message in CreateUnderCell %v", err)
	}

	return nil
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

	return &msg, err
}
