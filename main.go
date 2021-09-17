package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
)

var startKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Выбрать город"),
	),
)
var weatherKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(fmt.Sprintf("Погода на сегодня %v", "\U000026C4")),
		tgbotapi.NewKeyboardButton(fmt.Sprintf("Погода на завтра %v", "\U00002614")),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Поменять город"),
	),
)

func main() {

	cities := make(map[int64]string, 50)
	token := "2028283945:AAHWh2ms5PICriwlu-fvHx8jEDU-cfx_mFI"
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		// находимся ли мы в процессе выбора города

		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста нажмите выбрать город, чтобы начать")
		switch update.Message.Text {
		case "/start":
			// Если пользователь раннее вводил город
			if _, ok := cities[update.Message.Chat.ID]; ok {
				msg.ReplyMarkup = weatherKeyboard
			} else {
				msg.ReplyMarkup = startKeyboard
			}

		case "Выбрать город":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Напишите название города")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		case fmt.Sprintf("Погода на сегодня %v", "\U000026C4"):
			// Получаем и выводим погоду
			city := cities[update.Message.Chat.ID]
			err := GetWeather(city)
			if err != nil {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка проверки погоды")
			}

			temp, _ := GetTemperature(city, 0)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, strings.Title(temp))
		case fmt.Sprintf("Погода на завтра %v", "\U00002614"):
			// Получаем и выводим погоду
			city := cities[update.Message.Chat.ID]
			err := GetWeather(city)
			if err != nil {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка проверки погоды")
			}

			temp, _ := GetTemperature(city, 1)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, strings.Title(temp))
		case "Поменять город":
			// Удаляем существующий городу по ключу
			delete(cities, update.Message.Chat.ID)
			msg.ReplyMarkup = startKeyboard
		default:
			if _, ok := cities[update.Message.Chat.ID]; ok {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Некорректная команда")
				msg.ReplyMarkup = weatherKeyboard
			} else {
				cities[update.Message.Chat.ID] = update.Message.Text
				city := cities[update.Message.Chat.ID]
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Вы выбрали город "+city)
				err := GetWeather(city)
				msg.ReplyMarkup = weatherKeyboard
				if err != nil {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Некорректный город")
					msg.ReplyMarkup = startKeyboard
				}
			}

		}

		bot.Send(msg)
	}
}
