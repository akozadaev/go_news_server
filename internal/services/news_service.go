package services

import (
	"context"
	"go_news_server/internal/models"
	"go_news_server/internal/repository"
)

// NewsServiceInterface - интерфейс для NewsService
type NewsServiceInterface interface {
	UpdateNews(ctx context.Context, news *models.News, categories []int64) error
	GetNewsList(ctx context.Context, limit, offset int) ([]models.News, error)
}

type NewsService struct {
	Repository *repository.NewsRepository
}

// Проверка, что NewsService реализует интерфейс
var _ NewsServiceInterface = (*NewsService)(nil)

func (s *NewsService) UpdateNews(ctx context.Context, news *models.News, categories []int64) error {
	return s.Repository.UpdateNews(ctx, news, categories)
}

func (s *NewsService) GetNewsList(ctx context.Context, limit, offset int) ([]models.News, error) {
	return s.Repository.GetNewsList(ctx, limit, offset)
}
