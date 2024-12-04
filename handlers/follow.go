package handlers

import (
	"fmt"
	"net/http"
	"newsfeed-rohnsrma/db"
	"strconv"
)

func FollowUser(w http.ResponseWriter, r *http.Request) {
	followerID, err1 := strconv.Atoi(r.URL.Query().Get("follower_id"))
	followeeID, err2 := strconv.Atoi(r.URL.Query().Get("followee_id"))
	if err1 != nil || err2 != nil {
		http.Error(w, "Invalid user IDs", http.StatusBadRequest)
		return
	}

	_, err := db.PG.Exec("INSERT INTO follows (follower_id, followee_id) VALUES ($1, $2)", followerID, followeeID)

	if err != nil {
		http.Error(w, "Error following user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	cacheKey := fmt.Sprintf("newsfeed:user:%d", followerID)
	_, err = db.Redis.Do("DEL", cacheKey)
	if err != nil {
		fmt.Println("Error deleting cache key:", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Followed successfully"))
}
