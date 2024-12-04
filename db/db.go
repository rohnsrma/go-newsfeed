package db

import (
	"database/sql"
	"log"

	"github.com/gomodule/redigo/redis" // Redis client
	_ "github.com/lib/pq"              // PostgreSQL driver
)

var (
	PG    *sql.DB
	Redis redis.Conn
)

func Init() {
	var err error

	PG, err = sql.Open("postgres", "user=rohan dbname=newsfeed_db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if err = PG.Ping(); err != nil {
		log.Fatal(err)
	}

	Redis, err = redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
	}
}

func Close() {
	PG.Close()
	Redis.Close()
}
