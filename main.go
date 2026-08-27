package main

import (
	"context"
	"log"
	"neurodronizm/cmd/collector"
	"neurodronizm/cmd/store"
)

func main() {
	ctx := context.Background()
	s, err := store.New("postgres://dronism:change_me@localhost:5432/dronism?sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	collector.Collector(ctx, s)
}
