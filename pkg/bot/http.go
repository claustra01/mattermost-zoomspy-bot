package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// Client methods

func (c *Client) PostMessage(channelID string, message string) error {
	body := CreatePostRequestBody{
		ChannelID: channelID,
		Message:   message,
		RootID:    nil,
		FileIDs:   []string{},
	}
	buf, err := marshalJSON(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	resp, err := c.do(http.MethodPost, "/api/v4/posts", buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error response from server: StatusCode %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	slog.Info("Response:", "Body", string(respBody))
	return nil
}

func (c *Client) GetUserChannels() ([]Channel, error) {
	teams, err := c.fetchTeams()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch teams: %w", err)
	}

	var channels []Channel
	for _, team := range teams {
		teamChannels, err := c.fetchTeamChannels(team.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch channels for team %s: %w", team.ID, err)
		}
		channels = append(channels, teamChannels...)
	}
	return channels, nil
}

func (c *Client) GetUnreadPosts() ([]ChannelUnread, error) {
	channels, err := c.GetUserChannels()
	if err != nil {
		return nil, err
	}

	cutoff := time.Now().Add(-time.Hour).UnixMilli()

	var unread []ChannelUnread
	for _, channel := range channels {
		member, err := c.fetchChannelMember(channel.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch membership for channel %s: %w", channel.ID, err)
		}

		posts, err := c.fetchPostsSince(channel.ID, member.LastViewedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch posts for channel %s: %w", channel.ID, err)
		}

		var recentPosts []Post
		for _, post := range posts {
			if post.CreateAt >= cutoff {
				recentPosts = append(recentPosts, post)
			}
		}

		if len(recentPosts) == 0 {
			continue
		}

		unread = append(unread, ChannelUnread{
			Channel: channel,
			Posts:   recentPosts,
		})
	}

	return unread, nil
}

func (c *Client) GetUnreadZoomPosts() ([]ChannelUnread, error) {
	unread, err := c.GetUnreadPosts()
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

func (c *Client) MarkChannelRead(channelID string) error {
	payload := map[string]string{
		"channel_id": channelID,
	}
	buf, err := marshalJSON(payload)
	if err != nil {
		return err
	}

	resp, err := c.do(http.MethodPost, "/api/v4/channels/members/me/view", buf)
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

// Internal helpers

func (c *Client) fetchTeams() ([]Team, error) {
	resp, err := c.do(http.MethodGet, "/api/v4/users/me/teams", nil)
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
	if c.teamID == "" {
		return nil, fmt.Errorf("team ID is required")
	}

	for _, team := range teams {
		if team.ID == c.teamID {
			return []Team{team}, nil
		}
	}

	return nil, fmt.Errorf("team %s not found in user teams", c.teamID)
}

func (c *Client) fetchTeamChannels(teamID string) ([]Channel, error) {
	resp, err := c.do(http.MethodGet, fmt.Sprintf("/api/v4/users/me/teams/%s/channels", teamID), nil)
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

func (c *Client) fetchChannelMember(channelID string) (ChannelMember, error) {
	resp, err := c.do(http.MethodGet, fmt.Sprintf("/api/v4/channels/%s/members/me", channelID), nil)
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

func (c *Client) fetchPostsSince(channelID string, since int64) ([]Post, error) {
	resp, err := c.do(http.MethodGet, fmt.Sprintf("/api/v4/channels/%s/posts?since=%d", channelID, since), nil)
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
