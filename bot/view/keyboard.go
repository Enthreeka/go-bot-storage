package view

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var StartKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Показать мои разделы", "all_cell")),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Создать раздел", "create_cell")),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Удалить раздел", "delete_cell")))

var AddCellKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить раздел", "add_cell")))

var MainMenuButtonData = tgbotapi.NewInlineKeyboardButtonData("Вернуться в управление", "back_main")

var CellButtonDataCreate = tgbotapi.NewInlineKeyboardButtonData("Создать тему", "create_under_cell")

var CellButtonDataDelete = tgbotapi.NewInlineKeyboardButtonData("Удалить тему", "delete_under_cell")

var NoOneRowsButtonData = tgbotapi.NewInlineKeyboardButtonData("Пока что ничего нет", "no_one_rows")

var AddDataButtonData = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Добавить данные", "add_data"))

var DeleteDataButtonData = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Удалить данные", "delete_data"))
