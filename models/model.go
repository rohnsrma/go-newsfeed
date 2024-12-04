package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type Post struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type Follow struct {
	FollowerID int `json:"follower_id"`
	FolloweeID int `json:"followee_id"`
}
