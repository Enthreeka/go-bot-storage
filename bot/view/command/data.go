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

func (d *dataView) CreateData(update *tgbotapi.Update, msg *tgbotapi.MessageConfig, cellData *string) error {
	msg.ParseMode = tgbotapi.ModeHTML

	underCellID, _ := model.FindIdName(*cellData)

	err := d.dataController.CreateData(update, &underCellID)
	if err != nil {
		d.log.Error("failed to create in [add_data] by [%s]: %v", update.Message.From.UserName, err)

		msg.Text = "Ошибка при добавлении данных!"
		_, err = d.bot.Send(msg)
		if err != nil {
			d.log.Error("failed to create data in CreateData %v", err)
		}
		return err
	}

	msg.Text = "Данные добавлены успешно"
	_, err = d.bot.Send(msg)
	if err != nil {
		d.log.Error("failed to send message in CreateData %v", err)
	}

	return nil
}

func (d *dataView) UpdateData(update *tgbotapi.Update, msg *tgbotapi.MessageConfig, cellData *string) error {
	msg.ParseMode = tgbotapi.ModeHTML

	underCellID, _ := model.FindIdName(*cellData)

	err := d.dataController.UpdateData(update, &underCellID)
	if err != nil {
		d.log.Error("failed to show data in [update_data] by [%s]: %v", update.Message.From.UserName, err)

		msg.Text = "Ошибка при изменении данных!"
		_, err = d.bot.Send(msg)
		if err != nil {
			d.log.Error("failed to update data in UpdateData %v", err)
		}
		return err
	}

	msg.Text = "Данные изменены успешно"
	_, err = d.bot.Send(msg)
	if err != nil {
		d.log.Error("failed to send message in UpdateData %v", err)
	}

	return nil
}

func (d *dataView) ShowData(update *tgbotapi.Update, cellData *string) error {
	userID := update.Message.Chat.ID
	underCellID, name := model.FindIdName(*cellData)

	msg := tgbotapi.NewMessage(userID, "")
	msg.ParseMode = tgbotapi.ModeHTML

	data, err := d.dataController.GetData(underCellID)
	if err != nil {
		d.log.Error("failed to show data in [update_data] by [%s]: %v", update.Message.From.UserName, err)
		return err
	}
	markup := tgbotapi.NewInlineKeyboardMarkup(view.UpdateDataButtonData,
		tgbotapi.NewInlineKeyboardRow(view.MainMenuButtonData),
	)
	msg.ReplyMarkup = &markup

	dataFile, file := model.IsFile(data.Describe)
	if file {
		fileID := tgbotapi.FileID(dataFile)

		msg := tgbotapi.NewDocument(userID, fileID)
		msg.ReplyMarkup = &markup
		msg.ParseMode = tgbotapi.ModeHTML

		msg.Caption = fmt.Sprintf("<b>%s</b>\n", name)

		_, err = d.bot.Send(msg)
		if err != nil {
			d.log.Error("error sending document: %v", err)
			return err
		}
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
