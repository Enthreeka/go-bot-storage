package main

import (
	"github.com/Enthreeka/go-bot-storage/bot/controller"
	"github.com/Enthreeka/go-bot-storage/bot/model"
	"github.com/Enthreeka/go-bot-storage/bot/view/callback"
	"github.com/Enthreeka/go-bot-storage/bot/view/command"
	"github.com/Enthreeka/go-bot-storage/logger"
	"github.com/Enthreeka/go-bot-storage/repository/sqlite"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func main() {
	err := godotenv.Load("config/bot.env")
	if err != nil {
		log.Fatal("Error loading bot.env file")
	}

	tokenTelegram := os.Getenv("TG_TOKEN")

	log := logger.New()

	sqLite, err := sqlite.New()
	if err != nil {
		log.Fatal("failed to connect SqLite: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(tokenTelegram)
	if err != nil {
		log.Fatal("failed to load token %v", err)
	}

	bot.Debug = true

	log.Info("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	userRepo := sqlite.UserRepository(sqLite)

	userController := controller.UserController(userRepo, log)

	status := make(map[int64]*model.Status)
	for update := range updates {
		if update.Message != nil {
			log.Info("[%s] %s", update.Message.From.UserName, update.Message.Text)

			userID := update.Message.Chat.ID
			msg := tgbotapi.NewMessage(userID, update.Message.Text)

			// Check user subscriber or not
			_, err := userController.CheckUser(&update)
			if err != nil {
				log.Error("failed to check or create user: %v", err)
			}

			//Initialization map
			if _, ok := status[userID]; !ok {
				status[userID] = &model.Status{
					Callback: make(map[string]bool),
				}
			}

			switch update.Message.Command() {
			case "start":
				command.BotSendStart(log, bot, &msg)
			case "info":
				command.BotSendInfo(log, bot, &msg)
			default:
				if userStatus, ok := status[userID]; ok {
					if userStatus.Callback["create_cell"] == true {
						userStatus.Callback["create_cell"] = false

						//cell := update.Message.Text

					}
				} else {
					command.BotSendDefault(log, bot, &msg)
				}
			}
		} else if update.CallbackQuery != nil {
			userID := update.CallbackQuery.Message.Chat.ID

			switch update.CallbackQuery.Data {
			case "create_cell":

				status[userID].Callback["create_cell"] = true
				callback.BotSendTextCell(log, bot, userID)
			case "delete_cell":

			}
		}

	}

}
