package callback

import (
	"database/sql"
	"github.com/Enthreeka/go-bot-storage/bot/controller"
	"github.com/Enthreeka/go-bot-storage/bot/model"
	"github.com/Enthreeka/go-bot-storage/bot/view"
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

func (d *dataView) ShowData(update *tgbotapi.Update) (int, error) {
	userID := update.CallbackQuery.Message.Chat.ID
	underCellID, name := model.FindIntStr(update.CallbackQuery.Data)

	msg := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    userID,
			MessageID: update.CallbackQuery.Message.MessageID,
		}, Text: name,
	}

	data, err := d.dataController.GetData(underCellID)
	if err != nil {
		if err == sql.ErrNoRows {
			//TODO проверка, если пусто то вставить метод создания заметки
			markup := tgbotapi.NewInlineKeyboardMarkup(view.AddDataButtonData)
			msg.ReplyMarkup = &markup

			_, err = d.bot.Send(msg)
			if err != nil {
				d.log.Error("error sending under cell keyboard: %v", err)
			}
		}
		return underCellID, err
	}

	markup := tgbotapi.NewInlineKeyboardMarkup(view.AddDataButtonData, view.DeleteDataButtonData)
	msg.ReplyMarkup = &markup

	dataFile, file := model.IsFile(data.Describe)
	if file {
		fileID := tgbotapi.FileID(dataFile)
		msg := tgbotapi.DocumentConfig{
			BaseFile: tgbotapi.BaseFile{
				BaseChat: tgbotapi.BaseChat{
					ChatID:      userID,
					ReplyMarkup: &markup,
				},
				File: fileID,
			},
			Caption: name}

		_, err = d.bot.Send(msg)
		if err != nil {
			d.log.Error("error sending document: %v", err)
		}
		return underCellID, nil
	}

	msg.Text = data.Describe
	_, err = d.bot.Send(msg)
	if err != nil {
		d.log.Error("error sending text: %v", err)
	}

	return underCellID, nil
}
