package models

type News struct {
	Id         int64   `json:"id"`
	Title      string  `json:"title"`
	Content    string  `json:"content"`
	Categories []int64 `json:"categories"`
}
