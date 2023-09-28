package command

import (
	"errors"
	"fmt"
	"github.com/Enthreeka/go-bot-storage/bot/controller"
	"github.com/Enthreeka/go-bot-storage/bot/model"
	"github.com/Enthreeka/go-bot-storage/bot/view"
	"github.com/Enthreeka/go-bot-storage/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"time"
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
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(view.UpdateDataButtonData),
		tgbotapi.NewInlineKeyboardRow(view.RemindDataButtonData),
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

func (d *dataView) RemindData(update *tgbotapi.Update, data *string) error {
	userID := update.Message.Chat.ID
	msg := tgbotapi.NewMessage(userID, "")

	// Checking on correct input date
	callTime, err := time.ParseInLocation("15:04 02.01.2006", update.Message.Text, time.Local)
	if err != nil {
		d.log.Error("invalid date format by [%s]: %s", update.Message.From.UserName, update.Message.Text)
		msg.Text = "Неправильно введенный формат данных!"
		_, err := d.bot.Send(msg)
		if err != nil {
			d.log.Error("failed to send message in RemindData")
		}
		return err
	}

	// Checking is low then time now
	if !callTime.After(time.Now()) {
		d.log.Error("the requested time is longer than the present [%s]: %s", update.Message.From.UserName, update.Message.Text)
		msg.Text = "Неправильно введенный формат данных!"
		_, err := d.bot.Send(msg)
		if err != nil {
			d.log.Error("failed to send message in RemindData")
		}
		return errors.New("callTime is greater than time.Now()")
	}

	// Checking is date more than 1 year
	if callTime.After(time.Now().AddDate(1, 0, 0)) {
		d.log.Error("the requested time is more than one year in the future [%s]: %s", update.Message.From.UserName, update.Message.Text)
		msg.Text = "Неправильно введенный формат данных! Дата не может быть более чем через год."
		_, err := d.bot.Send(msg)
		if err != nil {
			d.log.Error("failed to send message in RemindData")
		}
		return errors.New("callTime is more than one year in the future")
	}

	duration := callTime.Sub(time.Now())
	waitTime := time.After(duration)
	for {
		select {
		case <-waitTime:
			underCellID, name := model.FindIdName(*data)

			data, err := d.dataController.GetData(underCellID)
			if err != nil {
				d.log.Error("failed to get data in reminder by [%s]: %v", update.Message.From.UserName, err)
				return err
			}

			dataFile, file := model.IsFile(data.Describe)
			if file {
				fileID := tgbotapi.FileID(dataFile)

				msg := tgbotapi.NewDocument(userID, fileID)
				msg.ParseMode = tgbotapi.ModeHTML

				msg.Caption = fmt.Sprintf("Спешу напомнить об:\n<b>%s</b>\n", name)

				_, err = d.bot.Send(msg)
				if err != nil {
					d.log.Error("error sending document: %v", err)
					return err
				}
				return nil
			}

			msg.ParseMode = tgbotapi.ModeHTML
			var builder strings.Builder
			builder.WriteString("Спешу напомнить об:\n")
			builder.WriteString("<b>")
			builder.WriteString(name)
			builder.WriteString("</b> ")
			builder.WriteString("\n\n")
			builder.WriteString(data.Describe)

			msg.Text = builder.String()

			_, err = d.bot.Send(msg)
			if err != nil {
				d.log.Error("failed to send message in RemindData")
				return err
			}
		}
	}

	return nil
}
