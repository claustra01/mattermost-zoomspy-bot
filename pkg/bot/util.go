package bot

import (
	"encoding/json"
)

type CreatePostRequestBody struct {
	ChannelID string   `json:"channel_id"`
	Message   string   `json:"message"`
	RootID    *string  `json:"root_id"`
	FileIDs   []string `json:"file_ids"`
}

func MarshalCreatePostReqBody(body CreatePostRequestBody) ([]byte, error) {
	raw, err := json.Marshal(body)
	if err != nil {
		return raw, err
	}
	return raw, nil
}
