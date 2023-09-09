package controller

import (
	"github.com/Enthreeka/go-bot-storage/logger"
	"github.com/Enthreeka/go-bot-storage/repository"
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
