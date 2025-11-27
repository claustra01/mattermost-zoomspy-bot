package main

import (
	"fmt"
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
				permalink := bot.BuildPostURL(baseUrl, post.ID)
				msg := fmt.Sprintf("[%s] %s", item.Channel.DisplayName, permalink)
				if err := bot.PostMessage(baseUrl, channelID, token, msg); err != nil {
					slog.Error("Error reposting zoom link", "channel_id", channelID, "error", err)
				}
			}
			if err := bot.MarkChannelRead(baseUrl, token, item.Channel.ID); err != nil {
				slog.Error("Error marking channel read", "channel_id", item.Channel.ID, "error", err)
			}
		}
	}

	job()

	// c := cron.New()
	// c.AddFunc("* * * * *", job)
	// c.Start()

	// select {}
}
