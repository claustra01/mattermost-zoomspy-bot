package bot

type Team struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Name        string `json:"name"`
}

type Channel struct {
	ID          string `json:"id"`
	TeamID      string `json:"team_id"`
	DisplayName string `json:"display_name"`
	Name        string `json:"name"`
	Type        string `json:"type"`
}

type ChannelMember struct {
	LastViewedAt int64 `json:"last_viewed_at"`
}

type Post struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	UserID    string `json:"user_id"`
	Message   string `json:"message"`
	RootID    string `json:"root_id"`
	ParentID  string `json:"parent_id"`
	CreateAt  int64  `json:"create_at"`
}

type PostsResponse struct {
	Order      []string        `json:"order"`
	Posts      map[string]Post `json:"posts"`
	NextPostID string          `json:"next_post_id"`
	PrevPostID string          `json:"prev_post_id"`
	HasMore    bool            `json:"has_more"`
}

type ChannelUnread struct {
	Channel Channel
	Posts   []Post
}

type CreatePostRequestBody struct {
	ChannelID string   `json:"channel_id"`
	Message   string   `json:"message"`
	RootID    *string  `json:"root_id"`
	FileIDs   []string `json:"file_ids"`
}
