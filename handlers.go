package main

import (
	"context"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	// `for {` means the loop is infinite until we manually stop it
	for {
		select {
		// stop looping if ctx is cancelled
		case <-ctx.Done():
			return
		// receive update from channel and then handle it
		case update := <-updates:
			handleUpdate(update)
		}
	}
}

func handleUpdate(update tgbotapi.Update) {
	if update.Message != nil {
		handleMessage(update.Message)
	}
}

func handleMessage(message *tgbotapi.Message) {
	user := message.From
	text := message.Text
	if message.Location != nil {
		lat := message.Location.Latitude
		lon := message.Location.Longitude
		_ = sendWeather(message.Chat.ID, lat, lon)
		return
	}

	if user == nil {
		return
	}

	log.Printf("%s wrote %s", user.FirstName, text)

	var err error

	if strings.HasPrefix(text, "/") {
		err = handleCommand(message.Chat.ID, text)
	} else {
		parts := strings.Fields(text)
		if len(parts) == 2 {
			latStr := strings.TrimSuffix(parts[0], ",")
			lonStr := parts[1]

			lat, err1 := strconv.ParseFloat(latStr, 64)
			lon, err2 := strconv.ParseFloat(lonStr, 64)
			if err1 != nil || err2 != nil {
				err = sendMessage(message.Chat.ID,
					"Некорректные координаты. Укажи так: 55.656201 37.587002 или отправь геопозицию ",
				)
				return
			}
			err = sendWeather(message.Chat.ID, lat, lon)
			return
		}

	}

	if err != nil {
		log.Printf("An error occured: %s", err.Error())
	}
}

// When we get a command, we react accordingly
func handleCommand(chatId int64, command string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return nil
	}

	cmd := parts[0]

	switch cmd {
	case "/start":
		return sendStart(chatId)
	case "/weather":
		if len(parts) < 3 {
			return sendMessage(chatId, "Пожалуйста укажите координаты в формате: /weather <широта> <долгота>")
		}

		latStr := strings.TrimSuffix(parts[1], ",")
		lonStr := parts[2]

		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			return sendMessage(chatId, "Некорректная широта")
		}
		lon, err := strconv.ParseFloat(lonStr, 64)
		if err != nil {
			return sendMessage(chatId, "Некорректная долгота")
		}
		return sendWeather(chatId, lat, lon)

	}

	return nil
}

func sendStart(chatId int64) error {
	locationButton := tgbotapi.NewKeyboardButtonLocation("Отправить геопозицию")
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(locationButton),
	)
	msg := tgbotapi.NewMessage(chatId, "Привет! Для отображения текущей погоды отправь координаты или свою геопозицию")
	msg.ReplyMarkup = keyboard
	_, err := bot.Send(msg)
	return err
}
