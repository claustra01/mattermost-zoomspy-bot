package main

import (
	"log/slog"
	"os"

	"github.com/claustra01/mattermost-zoomspy-bot/pkg/bot"
)

func main() {
	baseUrl := os.Getenv("MATTERMOST_BASE_URL")
	if baseUrl == "" {
		panic("MATTERMOST_BASE_URL is not set")
	}

	channelID := os.Getenv("MATTERMOST_CHANNEL_ID")
	if channelID == "" {
		panic("MATTERMOST_CHANNEL_ID is not set")
	}

	token := os.Getenv("MATTERMOST_TOKEN")
	if token == "" {
		panic("MATTERMOST_TOKEN is not set")
	}

	job := func() {
		err := bot.PostMessage(baseUrl, channelID, token)
		if err != nil {
			slog.Error("Error posting message:", "ERROR", err)
			return
		}
	}

	job()

	// c := cron.New()
	// c.AddFunc("* * * * *", job)
	// c.Start()

	// select {}
}
