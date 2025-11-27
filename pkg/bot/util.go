package bot

import (
	"encoding/json"
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

func MarshalCreatePostReqBody(body CreatePostRequestBody) ([]byte, error) {
	raw, err := json.Marshal(body)
	if err != nil {
		return raw, err
	}
	return raw, nil
}

func HasZoomURL(message string) bool {
	return zoomURLRegex.MatchString(message)
}

func BuildPostURL(baseUrl string, postID string) string {
	return strings.TrimSuffix(baseUrl, "/") + "/pl/" + postID
}
