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

	userRepo := sqlite.NewUserRepository(sqLite)
	cellRepo := sqlite.NewCellRepository(sqLite)

	userController := controller.NewUserController(userRepo, log)
	cellController := controller.NewCellController(cellRepo, log)

	cellViewCommand := command.NewCellView(cellController, bot, log)
	cellViewCallback := callback.NewCellView(cellController, bot, log)

	commandMail := command.NewCommandMail(cellController, bot, log)
	callbackMail := callback.NewCallbackMail(cellController, bot, log)

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
				commandMail.BotSendStart(&msg)
			case "info":
				commandMail.BotSendInfo(&msg)
			default:
				if userStatus, ok := status[userID]; ok {
					if userStatus.Callback["create_cell"] == true {
						userStatus.Callback["create_cell"] = false

						cellViewCommand.CreateCell(&update, &msg)

					}
				} else {
					commandMail.BotSendDefault(&msg)
				}
			}
		} else if update.CallbackQuery != nil {
			userID := update.CallbackQuery.Message.Chat.ID

			switch update.CallbackQuery.Data {
			case "create_cell":
				status[userID].Callback["create_cell"] = true
				callbackMail.BotSendTextCell(userID)
			case "delete_cell":
				status[userID].Callback["delete_cell"] = true
			case "all_cell":
				cellViewCallback.ShowCell(&update)
			}
		}

	}

}
