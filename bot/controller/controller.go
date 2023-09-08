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
	DeleteCell(name string) error
	GetCell(id int64) ([]model.Cell, error)
}
