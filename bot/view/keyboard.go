package view

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
	"unicode/utf8"
)

var AddCellKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить раздел", "add_cell")))

var CreateCellButtonData = tgbotapi.NewInlineKeyboardButtonData("Создать раздел", "create_cell")

var DeleteCellButtonData = tgbotapi.NewInlineKeyboardButtonData("Удалить раздел", "delete_cell")

var MainMenuButtonData = tgbotapi.NewInlineKeyboardButtonData("Вернуться в управление", "back_main")

var CellButtonDataCreate = tgbotapi.NewInlineKeyboardButtonData("Создать тему", "create_under_cell")

var CellButtonDataDelete = tgbotapi.NewInlineKeyboardButtonData("Удалить тему", "delete_under_cell")

var NoOneRowsButtonData = tgbotapi.NewInlineKeyboardButtonData("Пока что ничего нет", "no_one_rows")

var AddDataButtonData = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Добавить данные", "add_data"))

var UpdateDataButtonData = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Изменить данные", "update_data"))

func KeyboardValidation(text string) bool {
	// Checking for count symbol in text
	if utf8.RuneCountInString(text) > 37 {
		fmt.Println(utf8.RuneCountInString(text))
		return false
	}
	// Checking for symbol in ASCII table
	validPattern := regexp.MustCompile(`^[а-яА-Яa-zA-Z0-9\s\-_!@#$%^&*()]+`)

	return validPattern.MatchString(text)

}
