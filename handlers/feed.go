// handlers/feed.go
package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"net/http"
	"newsfeed-rohnsrma/db"
	"newsfeed-rohnsrma/models"
	"strconv"
)

func GetNewsFeed(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	cacheKey := fmt.Sprintf("newsfeed:user:%d", userID)

	cachedFeed, err := redis.Bytes(db.Redis.Do("GET", cacheKey))
	if err == nil {
		// Cache hit
		var posts []models.Post
		if err := json.Unmarshal(cachedFeed, &posts); err != nil {
			http.Error(w, "Error parsing cached data", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(posts)
		return
	}

	posts, err := fetchNewsFeedFromDB(userID)
	if err != nil {
		http.Error(w, "Error fetching news feed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	feedData, err := json.Marshal(posts)
	if err == nil {
		_, err = db.Redis.Do("SETEX", cacheKey, 60, feedData)
		if err != nil {
			fmt.Println("Error caching news feed:", err)
		}
	}

	json.NewEncoder(w).Encode(posts)
}

func fetchNewsFeedFromDB(userID int) ([]models.Post, error) {
	query := `
    SELECT posts.id, posts.user_id, posts.content, posts.created_at
    FROM posts
    JOIN follows ON posts.user_id = follows.followee_id
    WHERE follows.follower_id = $1
    ORDER BY posts.created_at DESC
    LIMIT 100
    `
	rows, err := db.PG.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
