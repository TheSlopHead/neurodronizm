package main

import (
	"context"
	"encoding/json"
	"log"
	"neurodronizm/internal/store"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Message struct {
	ID        int64  `json:"id"`
	Text      string `json:"text"`
	Date      int    `json:"date"`
	Forwarded string `json:"forwarded_from"`
}

type MessageExport struct {
	Message []Message `json:"messages"`
}

func main() {
	godotenv.Load()
	var export MessageExport
	ctx := context.Background()

	database_url := os.Getenv("DATABASE_URL")
	s, err := store.New(database_url)
	if err != nil {
		log.Printf("Cannot validate database: %v", err)
	}

	data, err := os.Open("data/result.json")
	if err != nil {
		log.Fatalf("Cannot read file: %v", err)
	}
	json.NewDecoder(data).Decode(&export)

	for _, msg := range export.Message {
		if msg.Forwarded != "" {
			continue
		}
		if msg.Text != "" {
			err := s.SavePost(ctx, msg.ID, msg.Text, time.Unix(int64(msg.Date), 0))
			if err != nil {
				log.Printf("Cannot save post: %v", err)
			}
		}

	}
}
