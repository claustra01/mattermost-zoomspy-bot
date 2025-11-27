package bot

import (
	"regexp"
	"strings"
)

type CreatePostRequestBody struct {
	ChannelID string   `json:"channel_id"`
	Message   string   `json:"message"`
	RootID    *string  `json:"root_id"`
	FileIDs   []string `json:"file_ids"`
}

var zoomURLRegex = regexp.MustCompile(`https?://[^\s]*zoom\.(us|com|gov)/[^\s]*`)

func HasZoomURL(message string) bool {
	return zoomURLRegex.MatchString(message)
}

func BuildPostURL(baseUrl string, postID string) string {
	return strings.TrimSuffix(baseUrl, "/") + "/pl/" + postID
}
