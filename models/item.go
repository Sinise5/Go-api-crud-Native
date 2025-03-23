package models

type Item struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}
