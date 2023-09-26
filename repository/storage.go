package repository

import "github.com/Enthreeka/go-bot-storage/bot/model"

type User interface {
	Create(user *model.User) (*model.User, error)
	GetByID(id int64) (*model.User, error)
}

type Cell interface {
	Create(cell *model.Cell) error
	DeleteByID(id int) error
	GetByUserID(id int64) ([]model.Cell, error)
}

type UnderCell interface {
	Create(cell *model.UnderCell) error
	DeleteByID(id int) error
	GetByCellID(userID int64, cellID int) ([]model.UnderCell, error)
}

type Data interface {
	Create(data *model.Data) error
	Update(data *model.Data) error
	GetByUnderCellID(underCellID int) (*model.Data, error)
	GetDataByName(dataName string, underCellID int) (*model.Data, error)
}
