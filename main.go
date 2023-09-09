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
	underCellRepo := sqlite.NewUnderCellRepository(sqLite)

	userController := controller.NewUserController(userRepo, log)
	cellController := controller.NewCellController(cellRepo, underCellRepo, log)

	cellViewCommand := command.NewCellView(cellController, bot, log)
	cellViewCallback := callback.NewCellView(cellController, bot, log)

	commandMail := command.NewCommandMail(bot, log)
	callbackMail := callback.NewCallbackMail(bot, log)

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

			//Initialization Callback map
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
			dataCommand := update.CallbackQuery.Data

			//cellID := make(map[int64]*int)

			// defines pre-defined buttons
			switch dataCommand {
			case "create_cell":
				status[userID].Callback["create_cell"] = true
				callbackMail.BotSendTextCell(userID)
			case "delete_cell":
				status[userID].Callback["delete_cell"] = true
			case "all_cell":
				cellViewCallback.ShowCell(&update)
			case "back_main":
				callbackMail.BotSendMainMenu(&update)
			case "create_under_cell":
			case "delete_under_cell":

			}

			// defines "cell_name_id" and "underCell_name_id" buttons
			if model.IsCell(dataCommand) {

				// TODO После нажатия на кнопки с разделами должны появится кнпоки с темами(либо же их нет)
				// TODO в кнопке с темами есть доп 3 кнопки -> вернуться в главное меню, создать тему, удалить тему
				//id, name := func(dataCommand string) (int, string) {
				//	re := regexp.MustCompile("[0-9]+")
				//	digits := re.FindAllString(dataCommand, -1)
				//	digitsStr := strings.Join(digits, "")
				//
				//	name := strings.Split(dataCommand, "_")
				//
				//	digitsInt, _ := strconv.Atoi(digitsStr)
				//
				//	return digitsInt, name[1]
				//}(dataCommand)

				//cellID[userID] = &id

				_, err := cellViewCallback.ShowUnderCell(&update)
				if err != nil {
					log.Error("%v", err)
				}

			} else if model.IsUnderCell(dataCommand) {

			}
		}

	}

}
