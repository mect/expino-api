package main

import "time"

// NewsItem contains the info for one newstitem
type NewsItem struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Date    time.Time `json:"date"`
}
