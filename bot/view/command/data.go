package command

import (
	"github.com/Enthreeka/go-bot-storage/bot/controller"
	"github.com/Enthreeka/go-bot-storage/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

func (d *dataView) CreateData(update *tgbotapi.Update, msg *tgbotapi.MessageConfig, underCellID *int) error {
	msg.ParseMode = tgbotapi.ModeHTML

	err := d.dataController.CreateData(update, underCellID)
	if err != nil {
		d.log.Error("failed to create [add_data] by [%s]: %v", update.Message.From.UserName, err)

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

func (d *dataView) UpdateData(update *tgbotapi.Update, msg *tgbotapi.MessageConfig, underCellID *int) error {
	msg.ParseMode = tgbotapi.ModeHTML

	err := d.dataController.UpdateData(update, underCellID)
	if err != nil {
		d.log.Error("failed to create [update_data] by [%s]: %v", update.Message.From.UserName, err)

		msg.Text = "Ошибка при изменении данных!"
		_, err = d.bot.Send(msg)
		if err != nil {
			d.log.Error("failed to update data in UpdateData %v", err)
		}
		return err
	}

	msg.Text = "Данные добавлены успешно"
	_, err = d.bot.Send(msg)
	if err != nil {
		d.log.Error("failed to send message in UpdateData %v", err)
	}

	return nil
}
