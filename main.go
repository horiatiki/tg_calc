package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <bot_token>")
		return
	}
	token := os.Args[1]
	botURL := apiURL + token + "/"

	offset := 0
	for {
		updates, err := getUpdates(botURL, offset)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		for _, upd := range updates {
			offset = upd.UpdateID + 1
			if upd.Message != nil && upd.Message.Text == "/start" {
				sendStartMessage(botURL, upd.Message.Chat.ID)
			} else if upd.CallbackQuery != nil {
				handleCallback(botURL, upd.CallbackQuery)
			}
		}
		time.Sleep(300 * time.Millisecond)
	}
}
