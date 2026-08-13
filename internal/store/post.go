package main

import (
	"time"
)

type Post struct {
	ID       int
	Text     string
	PostedAt time.Time
}

type PostStore struct {
	post   map[int]*Post
	nextID int
}

func CreatePost(store *PostStore) *Post {
	id := store.nextID
	store.nextID++
	time := time.Now()
	newPost := &Post{
		ID:       id,
		Text:     "lolkek",
		PostedAt: time,
	}
	if store.post == nil {
		store.post = make(map[int]*Post)
	}
	store.post[id] = newPost

	return newPost
}

func GetPost(*Post) {

}
