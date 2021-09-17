package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"pkg.re/essentialkaos/translit.v2"
)

var startKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Выбрать город"),
	),
)
var weatherKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Погода на сегодня"),
		tgbotapi.NewKeyboardButton("Погода на завтра"),

	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Поменять город"),
	),
)

func main() {

	// находимся ли мы в процессе выбора города
	var cityChoosen bool
	var city string

	bot, err := tgbotapi.NewBotAPI("2028283945:AAHWh2ms5PICriwlu-fvHx8jEDU-cfx_mFI")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста нажмите выбрать город, чтобы начать")

		switch update.Message.Text {
		case "/start":
			msg.ReplyMarkup = startKeyboard
		case "Выбрать город":
			cityChoosen = true
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Напишите название города")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		case "Погода на сегодня":
			temp, _ := GetTemperature(city)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Температура равна:  %f", temp))
		case "Поменять город":
			cityChoosen = true
			msg.ReplyMarkup = startKeyboard
		default:
			if cityChoosen {
				city = update.Message.Text
				city = translit.EncodeToISO9A(city)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Вы выбрали город " + city)
				err := GetWeather(city)
				msg.ReplyMarkup = weatherKeyboard
				if err != nil {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Некорректный город")
					msg.ReplyMarkup = startKeyboard
				}

				cityChoosen = false
			} else {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Некорректная команда")
			}

		}


		bot.Send(msg)
	}
}