package main

import "encoding/json"

func makeInlineKeyboard() string {
	rows := [][]string{
		{"7", "8", "9", "/"},
		{"4", "5", "6", "×"},
		{"1", "2", "3", "-"},
		{"0", ".", "=", "+"},
		{"C", "←", "√", "^2"},
	}
	inlineKeyboard := [][]map[string]string{}
	for _, row := range rows {
		buttons := []map[string]string{}
		for _, btn := range row {
			buttons = append(buttons, map[string]string{
				"text":          btn,
				"callback_data": btn,
			})
		}
		inlineKeyboard = append(inlineKeyboard, buttons)
	}

	markup := map[string]interface{}{
		"inline_keyboard": inlineKeyboard,
	}
	b, _ := json.Marshal(markup)
	return string(b)
}
