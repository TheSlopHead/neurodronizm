package main

import (
	"context"
	"log"
	"neurodronizm/cmd/collector"
	"neurodronizm/internal/generator"
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
	s, err := store.New(database_url)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	gen, _ := generator.New(context.Background(), os.Getenv("GEMINI_API_KEY"))
	ctx := context.Background()

	collector.Collector(ctx, s, gen)
}
