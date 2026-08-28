package main

import (
	"context"
	"log"
	"neurodronizm/internal/collector"
	"neurodronizm/internal/store"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	database_url := os.Getenv("DATABASE_URL")
	ctx := context.Background()
	s, err := store.New(database_url)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	collector.Collector(ctx, s)
}
