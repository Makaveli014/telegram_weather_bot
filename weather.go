package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type WeatherClient struct {
	apiKey string
	http   *http.Client
}

type WeatherResponse struct {
	Name    string `json:"name"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
	} `json:"main"`
}

func NewWeatherClient(apiKey string) *WeatherClient {
	return &WeatherClient{
		apiKey: apiKey,
		http:   &http.Client{Timeout: 3 * time.Second},
	}
}
func (w *WeatherClient) GetWeather(lat, lon float64) (string, error) {
	var data WeatherResponse
	base, _ := url.Parse("https://api.openweathermap.org/data/2.5/weather")
	params := url.Values{}
	params.Set("lat", fmt.Sprintf("%f", lat))
	params.Set("lon", fmt.Sprintf("%f", lon))
	params.Set("appid", w.apiKey)
	params.Set("units", "metric")
	params.Set("lang", "ru")
	base.RawQuery = params.Encode()
	finalURL := base.String()

	resp, err := w.http.Get(finalURL)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}
	desc := ""
	if len(data.Weather) > 0 {
		desc = data.Weather[0].Description
	}
	text := fmt.Sprintf(
		"%.0f°C, ощущается как %.0f°C, %s",
		data.Main.Temp,
		data.Main.FeelsLike,
		desc,
	)

	return text, nil
}

func sendWeather(chatId int64, lat, lon float64) error {
	action := tgbotapi.NewChatAction(chatId, tgbotapi.ChatTyping)
	bot.Send(action)
	text, err := weatherClient.GetWeather(lat, lon)
	if err != nil {
		return sendMessage(chatId, "Ошибка при получении погоды")
	}

	msg := tgbotapi.NewMessage(chatId, text)
	_, err = bot.Send(msg)
	return err
}

func sendMessage(chatId int64, message string) error {
	msg := tgbotapi.NewMessage(chatId, message)
	_, err := bot.Send(msg)
	return err
}
