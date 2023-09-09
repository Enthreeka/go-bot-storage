package controller

import (
	"github.com/Enthreeka/go-bot-storage/bot/model"
	"github.com/Enthreeka/go-bot-storage/logger"
	"github.com/Enthreeka/go-bot-storage/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type cellController struct {
	cellRepo      repository.Cell
	underCellRepo repository.UnderCell

	log *logger.Logger
}

func NewCellController(cellRepo repository.Cell, underCellRepo repository.UnderCell, log *logger.Logger) Cell {
	return &cellController{
		cellRepo:      cellRepo,
		underCellRepo: underCellRepo,
		log:           log,
	}
}

func (c *cellController) CreateCell(update *tgbotapi.Update) error {
	cell := &model.Cell{
		Name:   update.Message.Text,
		UserID: update.Message.Chat.ID,
	}

	err := c.cellRepo.Create(cell)
	if err != nil {
		return err
	}

	c.log.Info("cell creation has been successfully completed")
	return nil
}

func (c *cellController) DeleteCell(name string) error {
	//TODO implement me
	panic("implement me")
}

func (c *cellController) GetCell(id int64) ([]model.Cell, error) {
	cells, err := c.cellRepo.GetByUserID(id)
	if err != nil {
		return nil, err
	}

	return cells, nil
}

func (c *cellController) CreateUnderCell(update *tgbotapi.Update, cellID *int) error {
	cell := &model.UnderCell{
		CellID: *cellID,
		Name:   update.Message.Text,
	}

	err := c.underCellRepo.Create(cell)
	if err != nil {
		return err
	}

	return nil
}

func (c *cellController) DeleteUnderCell(name string) error {
	//TODO implement me
	panic("implement me")
}

func (c *cellController) GetUnderCell(userID int64, cellID int) ([]model.UnderCell, error) {
	underCell, err := c.underCellRepo.GetByCellID(userID, cellID)
	if err != nil {
		return nil, err
	}

	return underCell, err
}
