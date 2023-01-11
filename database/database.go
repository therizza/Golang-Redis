package database

import (
	"log"
	"os"

	redis "github.com/go-redis/redis/v9"
	"github.com/joho/godotenv"
)

var DB *redis.Client

func DbConnect() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	addr := os.Getenv("Addr")
	password := os.Getenv("Password")

	DB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

}
