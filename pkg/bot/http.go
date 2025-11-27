package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

func PostMessage(baseUrl string, channelID string, token string, message string) error {
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

	url := fmt.Sprintf("%s/api/v4/posts", trimBaseURL(baseUrl))
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

func GetUserChannels(baseUrl string, token string) ([]Channel, error) {
	teams, err := fetchTeams(baseUrl, token)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch teams: %w", err)
	}

	var channels []Channel
	for _, team := range teams {
		teamChannels, err := fetchTeamChannels(baseUrl, token, team.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch channels for team %s: %w", team.ID, err)
		}
		channels = append(channels, teamChannels...)
	}
	return channels, nil
}

func GetUnreadPosts(baseUrl string, token string) ([]ChannelUnread, error) {
	channels, err := GetUserChannels(baseUrl, token)
	if err != nil {
		return nil, err
	}

	var unread []ChannelUnread
	for _, channel := range channels {
		member, err := fetchChannelMember(baseUrl, token, channel.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch membership for channel %s: %w", channel.ID, err)
		}

		posts, err := fetchPostsSince(baseUrl, token, channel.ID, member.LastViewedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch posts for channel %s: %w", channel.ID, err)
		}

		if len(posts) == 0 {
			continue
		}

		unread = append(unread, ChannelUnread{
			Channel: channel,
			Posts:   posts,
		})
	}

	return unread, nil
}

func GetUnreadZoomPosts(baseUrl string, token string) ([]ChannelUnread, error) {
	unread, err := GetUnreadPosts(baseUrl, token)
	if err != nil {
		return nil, err
	}

	var filtered []ChannelUnread
	for _, item := range unread {
		var posts []Post
		for _, post := range item.Posts {
			if HasZoomURL(post.Message) {
				posts = append(posts, post)
			}
		}
		if len(posts) == 0 {
			continue
		}
		filtered = append(filtered, ChannelUnread{
			Channel: item.Channel,
			Posts:   posts,
		})
	}

	return filtered, nil
}

func MarkChannelRead(baseUrl string, token string, channelID string) error {
	payload := map[string]string{
		"channel_id": channelID,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal mark read payload: %w", err)
	}

	url := fmt.Sprintf("%s/api/v4/channels/members/me/view", trimBaseURL(baseUrl))
	resp, err := makeRequest(http.MethodPost, url, token, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to mark channel read: status %d body %s", resp.StatusCode, string(body))
	}
	return nil
}

func fetchTeams(baseUrl string, token string) ([]Team, error) {
	url := fmt.Sprintf("%s/api/v4/users/me/teams", trimBaseURL(baseUrl))
	resp, err := makeRequest("GET", url, token, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch teams: status %d body %s", resp.StatusCode, string(data))
	}

	var teams []Team
	if err := json.NewDecoder(resp.Body).Decode(&teams); err != nil {
		return nil, fmt.Errorf("failed to decode teams response: %w", err)
	}
	return teams, nil
}

func fetchTeamChannels(baseUrl string, token string, teamID string) ([]Channel, error) {
	url := fmt.Sprintf("%s/api/v4/users/me/teams/%s/channels", trimBaseURL(baseUrl), teamID)
	resp, err := makeRequest("GET", url, token, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch channels: status %d body %s", resp.StatusCode, string(data))
	}

	var channels []Channel
	if err := json.NewDecoder(resp.Body).Decode(&channels); err != nil {
		return nil, fmt.Errorf("failed to decode channels response: %w", err)
	}
	return channels, nil
}

func fetchChannelMember(baseUrl string, token string, channelID string) (ChannelMember, error) {
	url := fmt.Sprintf("%s/api/v4/channels/%s/members/me", trimBaseURL(baseUrl), channelID)
	resp, err := makeRequest("GET", url, token, nil)
	if err != nil {
		return ChannelMember{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return ChannelMember{}, fmt.Errorf("failed to fetch channel member: status %d body %s", resp.StatusCode, string(data))
	}

	var member ChannelMember
	if err := json.NewDecoder(resp.Body).Decode(&member); err != nil {
		return ChannelMember{}, fmt.Errorf("failed to decode channel member response: %w", err)
	}
	return member, nil
}

func fetchPostsSince(baseUrl string, token string, channelID string, since int64) ([]Post, error) {
	url := fmt.Sprintf("%s/api/v4/channels/%s/posts?since=%d", trimBaseURL(baseUrl), channelID, since)
	resp, err := makeRequest("GET", url, token, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch posts: status %d body %s", resp.StatusCode, string(data))
	}

	var postsResp PostsResponse
	if err := json.NewDecoder(resp.Body).Decode(&postsResp); err != nil {
		return nil, fmt.Errorf("failed to decode posts response: %w", err)
	}

	orderedPosts := make([]Post, 0, len(postsResp.Order))
	for _, id := range postsResp.Order {
		post, ok := postsResp.Posts[id]
		if !ok {
			continue
		}
		if post.CreateAt <= since {
			continue
		}
		orderedPosts = append(orderedPosts, post)
	}

	return orderedPosts, nil
}

func makeRequest(method string, url string, token string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	return resp, nil
}

func trimBaseURL(baseUrl string) string {
	return strings.TrimSuffix(baseUrl, "/")
}
