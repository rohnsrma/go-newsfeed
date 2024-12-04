package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"newsfeed-rohnsrma/db"
	"newsfeed-rohnsrma/models"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&post); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := db.PG.QueryRow("INSERT INTO posts (user_id, content) VALUES ($1, $2) RETURNING id, created_at",
		post.UserID, post.Content).Scan(&post.ID, &post.CreatedAt)
	if err != nil {
		http.Error(w, "Error creating post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	go invalidateFollowersCache(post.UserID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func invalidateFollowersCache(userID int) {
	rows, err := db.PG.Query("SELECT follower_id FROM follows WHERE followee_id = $1", userID)
	if err != nil {
		fmt.Println("Error fetching followers:", err)
		return
	}
	defer rows.Close()

	var followerID int
	for rows.Next() {
		err := rows.Scan(&followerID)
		if err != nil {
			fmt.Println("Error scanning follower ID:", err)
			continue
		}
		cacheKey := fmt.Sprintf("newsfeed:user:%d", followerID)
		_, err = db.Redis.Do("DEL", cacheKey)
		if err != nil {
			fmt.Println("Error deleting cache key:", err)
		}
	}
}
