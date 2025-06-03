package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const apiURL = "https://api.telegram.org/bot"

func getUpdates(botURL string, offset int) ([]Update, error) {
	resp, err := http.Get(fmt.Sprintf("%sgetUpdates?timeout=30&offset=%d", botURL, offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res Response
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil || !res.OK {
		return nil, fmt.Errorf("failed to get updates or bad response")
	}
	return res.Result, nil
}

func sendStartMessage(botURL string, chatID int64) {
	text := "Калькулятор\n\n"
	markup := makeInlineKeyboard()
	sendMessage(botURL, chatID, text, markup)
}

func sendMessage(botURL string, chatID int64, text, replyMarkup string) {
	data := url.Values{}
	data.Set("chat_id", strconv.FormatInt(chatID, 10))
	data.Set("text", text)
	data.Set("reply_markup", replyMarkup)
	data.Set("parse_mode", "Markdown")

	http.PostForm(botURL+"sendMessage", data)
}

func editMessageText(botURL string, chatID int64, messageID int, text, replyMarkup string) {
	data := url.Values{}
	data.Set("chat_id", strconv.FormatInt(chatID, 10))
	data.Set("message_id", strconv.Itoa(messageID))
	data.Set("text", text)
	data.Set("reply_markup", replyMarkup)
	data.Set("parse_mode", "Markdown")

	http.PostForm(botURL+"editMessageText", data)
}

func handleCallback(botURL string, cq *CallbackQuery) {
	chatID := cq.Message.Chat.ID
	msgID := cq.Message.MessageID
	data := cq.Data

	lines := strings.Split(cq.Message.Text, "\n")
	expr := ""
	if len(lines) > 1 {
		expr = lines[len(lines)-1]
	}

	switch data {
	case "C":
		editMessageText(botURL, chatID, msgID, "Калькулятор\n\n", makeInlineKeyboard())
	case "=":
		if !validExpression(expr) {
			editMessageText(botURL, chatID, msgID, "Калькулятор\n\nОшибка: неверное выражение", makeInlineKeyboard())
			return
		}
		res, err := evalExpression(expr)
		if err != nil {
			editMessageText(botURL, chatID, msgID, "Калькулятор\n\nОшибка: "+err.Error(), makeInlineKeyboard())
			return
		}
		editMessageText(botURL, chatID, msgID, fmt.Sprintf("Калькулятор\n\n%s = %v", expr, res), makeInlineKeyboard())
	case "√":
		if !validExpression(expr) {
			editMessageText(botURL, chatID, msgID, "Калькулятор\n\nОшибка: неверное выражение", makeInlineKeyboard())
			return
		}
		res, err := evalExpression(expr)
		if err != nil || res < 0 {
			editMessageText(botURL, chatID, msgID, "Калькулятор\n\nОшибка: неверное число для корня", makeInlineKeyboard())
			return
		}
		editMessageText(botURL, chatID, msgID, fmt.Sprintf("Калькулятор\n\n√(%s) = %v", expr, sqrt(res)), makeInlineKeyboard())
	case "^2":
		if !validExpression(expr) {
			editMessageText(botURL, chatID, msgID, "Калькулятор\n\nОшибка: неверное выражение", makeInlineKeyboard())
			return
		}
		res, err := evalExpression(expr)
		if err != nil {
			editMessageText(botURL, chatID, msgID, "Калькулятор\n\nОшибка: "+err.Error(), makeInlineKeyboard())
			return
		}
		editMessageText(botURL, chatID, msgID, fmt.Sprintf("Калькулятор\n\n(%s)^2 = %v", expr, res*res), makeInlineKeyboard())
	case "←":
		if len(expr) > 0 {
			expr = expr[:len(expr)-1]
		}
		editMessageText(botURL, chatID, msgID, "Калькулятор\n\n"+expr, makeInlineKeyboard())
	default:
		if len(lines) > 1 && strings.HasPrefix(lines[1], "Ошибка:") {
			// После ошибки разрешаем только C и ←
			return
		}
		if expr == "" || strings.Contains(expr, "=") {
			expr = ""
		}
		if isOperator(data) && (expr == "" || isOperator(string(expr[len(expr)-1]))) {
			// Не добавляем два оператора подряд
			editMessageText(botURL, chatID, msgID, "Калькулятор\n\n"+expr, makeInlineKeyboard())
			return
		}
		expr += data
		editMessageText(botURL, chatID, msgID, "Калькулятор\n\n"+expr, makeInlineKeyboard())
	}
}
