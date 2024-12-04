package main

import (
	"log"
	"net/http"
	"newsfeed-rohnsrma/db"
	"newsfeed-rohnsrma/handlers"
)

func main() {
	db.Init()
	defer db.Close()

	http.HandleFunc("/follow", handlers.FollowUser)
	http.HandleFunc("/post", handlers.CreatePost)
	http.HandleFunc("/feed", handlers.GetNewsFeed)

	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
