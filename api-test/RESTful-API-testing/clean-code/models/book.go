package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title        string `json:"title" form:"title"`
	Author       string `json:"author" form:"author"`
	Published_at string `json:"publishedAt" form:"publishedAt"`
}

type OutputBook struct {
	Title        string
	Author       string
	Published_at string
}
