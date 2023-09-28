package callback

import (
	"database/sql"
	"fmt"
	"github.com/Enthreeka/go-bot-storage/bot/controller"
	"github.com/Enthreeka/go-bot-storage/bot/model"
	"github.com/Enthreeka/go-bot-storage/bot/view"
	"github.com/Enthreeka/go-bot-storage/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type dataView struct {
	dataController controller.Data

	bot *tgbotapi.BotAPI
	log *logger.Logger
}

func NewDataView(dataController controller.Data, bot *tgbotapi.BotAPI, log *logger.Logger) *dataView {
	return &dataView{
		dataController: dataController,
		bot:            bot,
		log:            log,
	}
}

func (d *dataView) ShowData(update *tgbotapi.Update) error {
	userID := update.CallbackQuery.Message.Chat.ID
	underCellID, name := model.FindIdName(update.CallbackQuery.Data)

	msg := tgbotapi.EditMessageTextConfig{
		ParseMode: tgbotapi.ModeHTML,
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    userID,
			MessageID: update.CallbackQuery.Message.MessageID,
		},
	}

	data, err := d.dataController.GetData(underCellID)
	if err != nil {
		if err == sql.ErrNoRows {
			markup := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(view.AddDataButtonData),
				tgbotapi.NewInlineKeyboardRow(view.MainMenuButtonData))
			msg.ReplyMarkup = &markup
			msg.Text = fmt.Sprintf("<b>%s</b>", name)

			_, err = d.bot.Send(msg)
			if err != nil {
				d.log.Error("error sending under cell keyboard: %v", err)
			}
			return nil
		}
		d.log.Error("failed to get data in [underCell_name_id] by [%s]: %v", update.CallbackQuery.Message.From.UserName, err)
		return err
	}

	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(view.UpdateDataButtonData),
		tgbotapi.NewInlineKeyboardRow(view.RemindDataButtonData),
		tgbotapi.NewInlineKeyboardRow(view.MainMenuButtonData))
	msg.ReplyMarkup = &markup

	dataFile, file := model.IsFile(data.Describe)
	if file {
		fileID := tgbotapi.FileID(dataFile)
		msg := tgbotapi.DocumentConfig{
			ParseMode: tgbotapi.ModeHTML,
			BaseFile: tgbotapi.BaseFile{
				BaseChat: tgbotapi.BaseChat{
					ChatID:      userID,
					ReplyMarkup: &markup,
				},
				File: fileID,
			},
		}

		msg.Caption = fmt.Sprintf("<b>%s</b>\n", name)

		_, err = d.bot.Send(msg)
		if err != nil {
			d.log.Error("error sending document: %v", err)
			return err
		}

		d.log.Info("[%s] received document", update.CallbackQuery.Message.From.UserName)
		return nil
	}

	var builder strings.Builder
	builder.WriteString("<b>")
	builder.WriteString(name)
	builder.WriteString("</b> ")
	builder.WriteString("\n\n")
	builder.WriteString(data.Describe)

	msg.Text = builder.String()
	_, err = d.bot.Send(msg)
	if err != nil {
		d.log.Error("error sending text: %v", err)
	}

	return nil
}
