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

func (d *dataView) ShowData(update *tgbotapi.Update) (int, error) {
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
			//TODO проверка, если пусто то вставить метод создания заметки
			//TODO изменить NewInlineKeyboardRow
			markup := tgbotapi.NewInlineKeyboardMarkup(view.AddDataButtonData,
				tgbotapi.NewInlineKeyboardRow(view.MainMenuButtonData))
			msg.ReplyMarkup = &markup
			msg.Text = fmt.Sprintf("<b>%s</b>", name)

			_, err = d.bot.Send(msg)
			if err != nil {
				d.log.Error("error sending under cell keyboard: %v", err)
			}
		}
		return underCellID, err
	}

	//TODO изменить структуру кнопок
	markup := tgbotapi.NewInlineKeyboardMarkup(view.UpdateDataButtonData,
		tgbotapi.NewInlineKeyboardRow(view.MainMenuButtonData),
	)
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
		}
		return underCellID, nil
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

	return underCellID, nil
}
