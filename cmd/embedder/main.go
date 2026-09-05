package main

import (
	"context"
	"fmt"
	"log"
	"neurodronizm/internal/generator"
	"neurodronizm/internal/store"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	ctx := context.Background()

	database_url := os.Getenv("DATABASE_URL")
	s, err := store.New(database_url)
	if err != nil {
		log.Fatalf("failed to validate database: %v", err)
	}
	client, err := generator.New(ctx, "ollama")

	posts, err := s.GetPostsWithoutEmbeddings(ctx)
	if err != nil {
		log.Fatalf("Get posts without embeddings error: %v", err)
	}

	for _, p := range posts {
		itemctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

		vector, err := client.GetEmbedding(itemctx, p.Text)
		if err != nil {
			log.Fatalf("get embedding error: %v", err)
		}
		err = s.SaveEmbedding(itemctx, p.ID, vector)
		if err != nil {
			log.Fatalf("save embedding error: %v", err)
		}
		fmt.Printf("Post %d vectorizing succesful!", p.ID)
		cancel()
		time.Sleep(200 * time.Millisecond)

	}

}
