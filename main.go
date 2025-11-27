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

	token := os.Getenv("MATTERMOST_TOKEN")
	if token == "" {
		panic("MATTERMOST_TOKEN is not set")
	}

	job := func() {
		channels, err := bot.GetUserChannels(baseUrl, token)
		if err != nil {
			slog.Error("Error fetching channels", "error", err)
			return
		}
		slog.Info("Fetched channels", "count", len(channels))

		unread, err := bot.GetUnreadZoomPosts(baseUrl, token)
		if err != nil {
			slog.Error("Error fetching zoom posts", "error", err)
			return
		}

		for _, item := range unread {
			slog.Info("Zoom posts found", "channel_id", item.Channel.ID, "channel", item.Channel.DisplayName, "count", len(item.Posts))
			for _, post := range item.Posts {
				slog.Info("Post", "channel_id", item.Channel.ID, "post_id", post.ID, "user_id", post.UserID, "created_at", post.CreateAt, "message", post.Message)
			}
		}
	}

	job()

	// c := cron.New()
	// c.AddFunc("* * * * *", job)
	// c.Start()

	// select {}
}
