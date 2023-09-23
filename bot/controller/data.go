package controller

import (
	"database/sql"
	"fmt"
	"github.com/Enthreeka/go-bot-storage/bot/model"
	"github.com/Enthreeka/go-bot-storage/logger"
	"github.com/Enthreeka/go-bot-storage/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type dataController struct {
	dataRepo repository.Data

	log *logger.Logger
}

func NewDataController(dataRepo repository.Data, log *logger.Logger) Data {
	return &dataController{
		dataRepo: dataRepo,
		log:      log,
	}
}

func (d *dataController) GetData(underCellID int) (*model.Data, error) {
	data, err := d.dataRepo.GetByUnderCellID(underCellID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return data, nil
}

func (d *dataController) CreateData(update *tgbotapi.Update, UnderCellID *int) error {
	data := &model.Data{
		UnderCellID: *UnderCellID,
	}

	if update.Message.Text != "" {
		data.Describe = update.Message.Text
	} else if update.Message.Document != nil {
		describe := fmt.Sprintf("file-%s", update.Message.Document.FileID)
		data.Describe = describe
	}

	err := d.dataRepo.Create(data)
	if err != nil {
		return err
	}

	return nil
}
