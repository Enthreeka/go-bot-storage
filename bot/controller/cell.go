package controller

import (
	"context"
	"github.com/Enthreeka/go-bot-storage/bot/model"
	"github.com/Enthreeka/go-bot-storage/logger"
	"github.com/Enthreeka/go-bot-storage/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type cellController struct {
	cellRepoSqlite repository.Cell
	underCellRepo  repository.UnderCell

	log *logger.Logger
}

func NewCellController(cellRepo repository.Cell, underCellRepo repository.UnderCell, log *logger.Logger) Cell {
	return &cellController{
		cellRepoSqlite: cellRepo,
		underCellRepo:  underCellRepo,
		log:            log,
	}
}

func (c *cellController) CreateCell(update *tgbotapi.Update) error {
	cell := &model.Cell{
		Name:   update.Message.Text,
		UserID: update.Message.Chat.ID,
	}

	err := c.cellRepoSqlite.Create(context.Background(), cell)
	if err != nil {
		return err
	}

	c.log.Info("[%s] cell creation has been successfully completed", update.Message.From.UserName)
	return nil
}

func (c *cellController) DeleteCell(id int) error {
	err := c.cellRepoSqlite.DeleteByID(context.Background(), id)
	if err != nil {
		return err
	}

	return nil
}

func (c *cellController) DeleteUnderCell(id int) error {
	err := c.underCellRepo.DeleteByID(context.Background(), id)
	if err != nil {
		return err
	}

	return nil
}

func (c *cellController) GetCell(id int64) ([]model.Cell, error) {
	cells, err := c.cellRepoSqlite.GetByUserID(context.Background(), id)
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

	err := c.underCellRepo.Create(context.Background(), cell)
	if err != nil {
		return err
	}

	return nil
}

func (c *cellController) GetUnderCell(userID int64, cellID int) ([]model.UnderCell, error) {
	underCell, err := c.underCellRepo.GetByCellID(context.Background(), userID, cellID)
	if err != nil {
		return nil, err
	}

	return underCell, err
}
