package main

type Update struct {
	UpdateID      int            `json:"update_id"`
	Message       *Message       `json:"message,omitempty"`
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text"`
}

type Chat struct {
	ID int64 `json:"id"`
}

type CallbackQuery struct {
	ID      string  `json:"id"`
	Message Message `json:"message"`
	Data    string  `json:"data"`
}

type Response struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}
