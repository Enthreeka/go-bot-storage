package controller

import (
	"github.com/Enthreeka/go-bot-storage/bot/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type User interface {
	CheckUser(update *tgbotapi.Update) (*model.User, error)
}

type Cell interface {
	CreateCell(update *tgbotapi.Update) error
	DeleteCell(id int) error
	GetCell(id int64) ([]model.Cell, error)
	UnderCell
}

type UnderCell interface {
	CreateUnderCell(update *tgbotapi.Update, cellID *int) error
	DeleteUnderCell(name string) error
	GetUnderCell(userID int64, cellID int) ([]model.UnderCell, error)
}

type Data interface {
	GetData(underCellID int) (*model.Data, error)
	CreateData(update *tgbotapi.Update, UnderCellID *int) error
}
