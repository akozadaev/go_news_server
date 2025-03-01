package services

import (
	"context"
	"go_news_server/internal/models"
	"go_news_server/internal/repository"
)

type NewsService struct {
	Repository *repository.NewsRepository
}

func (s *NewsService) UpdateNews(ctx context.Context, news *models.News) error {
	return s.Repository.UpdateNews(ctx, news)
}

func (s *NewsService) GetNewsList(ctx context.Context, limit, offset int) ([]models.News, error) {
	return s.Repository.GetNewsList(ctx, limit, offset)
}

func (s *NewsService) AddNews(ctx context.Context, news *models.News) error {
	return s.Repository.AddNews(ctx, news)
}
