package services

import (
	"context"
	"go_news_server/internal/models"
	"go_news_server/internal/repository"
)

type NewsService interface {
	GetNews(ctx context.Context) ([]models.News, error)
	UpdateNews(ctx context.Context, news *models.News) error
}

type newsService struct {
	repo repository.NewsRepository
}

func NewNewsService(repo repository.NewsRepository) NewsService {
	return &newsService{repo: repo}
}

func (s *newsService) GetNews(ctx context.Context) ([]models.News, error) {
	return s.repo.GetNews(ctx)
}

func (s *newsService) UpdateNews(ctx context.Context, news *models.News) error {
	return s.repo.UpdateNews(ctx, news)
}
