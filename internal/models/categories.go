package models

import (
	"time"
)

// Category представляет категорию новостей
type Category struct {
	Id          int64     `json:"id" db:"Id"`
	Name        string    `json:"name" db:"Name"`
	Description string    `json:"description" db:"Description"`
	CreatedAt   time.Time `json:"created_at" db:"CreatedAt"`
	UpdatedAt   time.Time `json:"updated_at" db:"UpdatedAt"`
}

// CategoryCreateRequest представляет запрос на создание категории
type CategoryCreateRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description"`
}

// CategoryUpdateRequest представляет запрос на обновление категории
type CategoryUpdateRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description"`
}

// CategoryResponse представляет ответ с категорией
type CategoryResponse struct {
	Success  bool      `json:"success"`
	Category *Category `json:"category,omitempty"`
	Message  string    `json:"message,omitempty"`
}

// CategoriesResponse представляет ответ со списком категорий
type CategoriesResponse struct {
	Success    bool       `json:"success"`
	Categories []Category `json:"categories"`
	Total      int64      `json:"total"`
}
