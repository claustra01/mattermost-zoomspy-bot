package bot

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

func PostMessage(baseUrl string, channelID string, token string) error {
	// create message
	message := "test message"

	// create request body
	body := CreatePostRequestBody{
		ChannelID: channelID,
		Message:   message,
		RootID:    nil,
		FileIDs:   []string{},
	}
	data, err := MarshalCreatePostReqBody(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	// create request
	url := fmt.Sprintf("%s/api/v4/posts", baseUrl)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error response from server: StatusCode %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	slog.Info("Response:", "URL", url, "Body", string(respBody))
	return nil
}
