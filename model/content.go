package model

type Content struct {
	ID       string `json:"id"`
	Sort_key string `json:"sort_key"`
	Title    string `json:"title"`
	Date     string `json:"date"`
	Category string `json:"category"`
	Type     string `json:"type"`
	Image    string `json:"image"`
	Creator  string `json:"creator"`
	Content  string `json:"content"`
	Active   string `json:"active"`
	Link     string `json:"link"`
}
