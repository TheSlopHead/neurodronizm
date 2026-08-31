package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

func main() {
	godotenv.Load()
	neuro_api := os.Getenv("GEMINI_API_KEY")
	database_url := os.Getenv("DATABASE_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  neuro_api,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatalf("Cannot create gemini client:%v", err)
	}

	db, err := sql.Open("pgx", database_url)
	if err != nil {
		log.Fatalf("Cannot open database: %v", err)
	}
	query := "SELECT tg_message_id, text, posted_at FROM posts ORDER BY tg_message_id ASC LIMIT 3; "
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Fatalf("Error or timeout: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int64
		var text string
		var postedAt time.Time

		err := rows.Scan(&id, &text, &postedAt)
		if err != nil {
			log.Fatalf("Cannot read string: %v", err)
		}
		prompt := "напиши пост на тему машин в стиле этих трех постов, постиронично и со стебом"

		parts := []*genai.Part{
			{Text: text + prompt},
		}
		result, err := client.Models.GenerateContent(ctx, "gemini-3.6-flash", []*genai.Content{{Parts: parts}}, nil)
		if err != nil {
			log.Fatalf("Cannot generate post: %v", err)
		}

		fmt.Println(result)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error with strings:%v", err)
	}
}
