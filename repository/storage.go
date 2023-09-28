package repository

import (
	"context"
	"github.com/Enthreeka/go-bot-storage/bot/model"
)

type User interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
}

type Cell interface {
	Create(ctx context.Context, cell *model.Cell) error
	DeleteByID(ctx context.Context, id int) error
	GetByUserID(ctx context.Context, id int64) ([]model.Cell, error)
}

type UnderCell interface {
	Create(ctx context.Context, cell *model.UnderCell) error
	DeleteByID(ctx context.Context, id int) error
	GetByCellID(ctx context.Context, userID int64, cellID int) ([]model.UnderCell, error)
}

type Data interface {
	Create(ctx context.Context, data *model.Data) error
	Update(ctx context.Context, data *model.Data) error
	GetByUnderCellID(ctx context.Context, underCellID int) (*model.Data, error)
	GetDataByName(ctx context.Context, dataName string, underCellID int) (*model.Data, error)
}
