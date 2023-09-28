package main

import (
	"context"
	"github.com/Enthreeka/go-bot-storage/bot/controller"
	"github.com/Enthreeka/go-bot-storage/bot/model"
	"github.com/Enthreeka/go-bot-storage/bot/view/callback"
	"github.com/Enthreeka/go-bot-storage/bot/view/command"
	"github.com/Enthreeka/go-bot-storage/config"
	"github.com/Enthreeka/go-bot-storage/logger"
	"github.com/Enthreeka/go-bot-storage/repository/redis"
	"github.com/Enthreeka/go-bot-storage/repository/sqlite"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log := logger.New()

	cfg, err := config.New()
	if err != nil {
		log.Fatal("failed to load config: %v", err)
	}

	sqLite, err := sqlite.New()
	if err != nil {
		log.Fatal("failed to connect SqLite: %v", err)
	}

	// Connect to Redis
	rds, err := redis.New(context.Background(), cfg)
	if err != nil {
		log.Error("redis is not working: %v", err)
	}
	defer rds.Close()

	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Fatal("failed to load token %v", err)
	}

	bot.Debug = false

	log.Info("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	userRepo := sqlite.NewUserRepository(sqLite)
	cellRepo := sqlite.NewCellRepositorySL(sqLite)
	underCellRepo := sqlite.NewUnderCellRepository(sqLite)
	dataRepo := sqlite.NewDataRepository(sqLite)

	userRepoRedis := redis.NewUserRepositoryRedis(rds)

	userController := controller.NewUserController(userRepo, userRepoRedis, log)
	cellController := controller.NewCellController(cellRepo, underCellRepo, log)
	dataController := controller.NewDataController(dataRepo, log)

	cellViewCommand := command.NewCellView(cellController, bot, log)
	dataViewCommand := command.NewDataView(dataController, bot, log)
	cellViewCallback := callback.NewCellView(cellController, bot, log)
	dataViewCallback := callback.NewDataView(dataController, bot, log)

	commandMail := command.NewCommandMail(bot, log)
	callbackMail := callback.NewCallbackMail(bot, log)

	status := make(map[int64]*model.Status)

	// Maps store update.CallbackQuery.Data
	cellData := make(map[int64]*string)
	underCellData := make(map[int64]*string)

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

			//Initialization Callback redis
			if _, ok := status[userID]; !ok {
				status[userID] = &model.Status{
					Callback: make(map[string]bool),
				}
			}

			switch update.Message.Command() {
			case "start":
				//Drop all user state
				if _, ok := status[userID]; ok {
					status[userID].Callback["delete_cell"] = false
					status[userID].Callback["delete_under_cell"] = false
					status[userID].Callback["update_data"] = false
					status[userID].Callback["remind_data"] = false
				}

				cellViewCommand.ShowCell(&update, &msg)
			case "info":
				commandMail.BotSendInfo(&msg)
			default:
				if userStatus, ok := status[userID]; ok {
					if userStatus.Callback["create_cell"] == true {
						userStatus.Callback["create_cell"] = false
						log.Info("[%s] Message-[create_cell]", update.Message.From.UserName)

						cellViewCommand.CreateCell(&update, &msg)
						cellViewCommand.ShowCell(&update, &msg)
					} else if userStatus.Callback["create_under_cell"] == true {
						userStatus.Callback["create_under_cell"] = false
						log.Info("[%s] Message-[create_under_cell]", update.Message.From.UserName)

						if cell, ok := cellData[userID]; ok {
							cellViewCommand.CreateUnderCell(&update, &msg, cell)
							cellViewCommand.ShowUnderCell(&update, cell)
						}
					} else if userStatus.Callback["add_data"] == true {
						userStatus.Callback["add_data"] = false
						log.Info("[%s] Message-[add_data]", update.Message.From.UserName)

						if underCellID, ok := underCellData[userID]; ok {
							dataViewCommand.CreateData(&update, &msg, underCellID)
							dataViewCommand.ShowData(&update, underCellID)
						}
					} else if userStatus.Callback["update_data"] == true {
						userStatus.Callback["update_data"] = false
						log.Info("[%s] Message-[update_data]", update.Message.From.UserName)

						if underCell, ok := underCellData[userID]; ok {
							dataViewCommand.UpdateData(&update, &msg, underCell)
							dataViewCommand.ShowData(&update, underCell)
						}
					} else if userStatus.Callback["remind_data"] == true {
						userStatus.Callback["remind_data"] = false
						log.Info("[%s] Message-[remind_data]", update.Message.From.UserName)

						if underCell, ok := underCellData[userID]; ok {
							go dataViewCommand.RemindData(&update, underCell)
						}
					} else {
						commandMail.BotSendDefault(&msg)
					}
				}
			}

		} else if update.CallbackQuery != nil {
			userID := update.CallbackQuery.Message.Chat.ID
			dataCommand := update.CallbackQuery.Data

			//Initialization Callback redis for case with restart bot
			if _, ok := status[userID]; !ok {
				status[userID] = &model.Status{
					Callback: make(map[string]bool),
				}
			}

			// defines pre-defined buttons
			switch dataCommand {
			case "create_cell":
				log.Info("[%s] CallbackQuery-[%s]", update.CallbackQuery.From.UserName, dataCommand)

				status[userID].Callback["create_cell"] = true
				callbackMail.BotSendTextCell(userID)
			case "delete_cell":
				log.Info("[%s] CallbackQuery-[%s]", update.CallbackQuery.From.UserName, dataCommand)

				status[userID].Callback["delete_cell"] = true
				callbackMail.BotSendTextDeleteCell(userID)
			case "back_main":
				log.Info("[%s] CallbackQuery-[%s]", update.CallbackQuery.From.UserName, dataCommand)

				cellViewCallback.ShowCell(&update)
			case "create_under_cell":
				log.Info("[%s] CallbackQuery-[%s]", update.CallbackQuery.From.UserName, dataCommand)

				status[userID].Callback["create_under_cell"] = true
				callbackMail.BotSendTextUnderCell(userID)
			case "delete_under_cell":
				log.Info("[%s] CallbackQuery-[%s]", update.CallbackQuery.From.UserName, dataCommand)

				status[userID].Callback["delete_under_cell"] = true
				callbackMail.BotSendTextDeleteUnderCell(userID)
			case "add_data":
				log.Info("[%s] CallbackQuery-[%s]", update.CallbackQuery.From.UserName, dataCommand)

				status[userID].Callback["add_data"] = true
				callbackMail.BotSendTextData(userID)
			case "update_data":
				log.Info("[%s] CallbackQuery-[%s]", update.CallbackQuery.From.UserName, dataCommand)

				status[userID].Callback["update_data"] = true
				callbackMail.BotSendTextUpdateData(userID)
			case "remind_data":
				log.Info("[%s] CallbackQuery-[%s]", update.CallbackQuery.From.UserName, dataCommand)

				status[userID].Callback["remind_data"] = true
				callbackMail.BotSendTextRemindData(userID)
			}

			// Defines "cell_name_id" , "underCell_name_id" , "delete_cell" , "delete_under_cell" buttons
			if model.IsCell(dataCommand) {
				// Checking for delete cell user state. It is "delete_cell" button
				if status[userID].Callback["delete_cell"] == true {
					status[userID].Callback["delete_cell"] = false
					log.Info("[%s] delete Cell - [%s]", update.CallbackQuery.From.UserName, dataCommand)

					msg, _ := cellViewCallback.DeleteCell(&update)
					//TODO посмотреть можно ли использоваь cellViewCallback, а не cellViewCommand
					cellViewCommand.ShowCell(&update, msg)
					// When user does not press "delete_cell" is executed display list cells. Is "cell_name_id" button
				} else {
					log.Info("[%s] open Cell - [%s]", update.CallbackQuery.From.UserName, dataCommand)
					err := cellViewCallback.ShowUnderCell(&update, "")
					if err != nil {
						log.Error("%v", err)
					}
					cellData[userID] = &update.CallbackQuery.Data
				}
			} else if model.IsUnderCell(dataCommand) {
				if status[userID].Callback["delete_under_cell"] == true {
					status[userID].Callback["delete_under_cell"] = false
					log.Info("[%s] delete UnderCell - [%s]", update.CallbackQuery.From.UserName, dataCommand)

					cellViewCallback.DeleteUnderCell(&update)
					if cell, ok := cellData[userID]; ok {
						cellViewCallback.ShowUnderCell(&update, *cell)
					}
					// When user does not press "delete_under_cell" is executed display list cells. Is "underCell_name_id" button
				} else {
					log.Info("[%s] open UnderCell - [%s]", update.CallbackQuery.From.UserName, dataCommand)
					err := dataViewCallback.ShowData(&update)
					if err != nil {
						log.Error("%v", err)
					}
					underCellData[userID] = &update.CallbackQuery.Data
				}
			}
		}
	}
}
