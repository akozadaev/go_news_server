package repository

import (
	"context"
	"errors"
	"go_news_server/internal/models"
	"gopkg.in/reform.v1"
)

// https://github.com/go-reform/reform/blob/main/querier_examples_test.go

type NewsRepository interface {
	GetNews(ctx context.Context) ([]models.News, error)
	UpdateNews(ctx context.Context, news *models.News) error
}

type newsRepository struct {
	db *reform.DB
}

func NewNewsRepository(db *reform.DB) NewsRepository {
	return &newsRepository{db: db}
}

func (r *newsRepository) GetNews(ctx context.Context) ([]models.News, error) {
	var news []models.News
	_, err := r.db.FindAllFrom(models.NewsTable, "ORDER BY Id", &news)
	if err != nil {
		return nil, err
	}
	return news, nil
}

func (r *newsRepository) UpdateNews(ctx context.Context, news *models.News) error {
	if news.ID == 0 {
		return errors.New("invalid ID")
	}
	err := r.db.Update(news)
	return err
}
