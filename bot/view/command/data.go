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

func (d *dataView) CreateData(update *tgbotapi.Update, msg *tgbotapi.MessageConfig, underCellID *int) {

	err := d.dataController.CreateData(update, underCellID)
	if err != nil {
		d.log.Error("%v", err)
	}

}
