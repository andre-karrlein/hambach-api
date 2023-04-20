package model

type Article struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Date     string `json:"date"`
	Content  string `json:"content"`
	Category string `json:"category"`
	Type     string `json:"type"`
	Image    string `json:"image"`
	Link     string `json:"link"`
}
