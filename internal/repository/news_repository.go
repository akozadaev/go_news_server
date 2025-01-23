package repository

import (
	"context"
	"errors"
	"fmt"
	"go_news_server/internal/models"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
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
	str, err := r.db.SelectAllFrom(models.NewsTable, "")
	for _, n := range str {
		row := n.(*models.News)
		str1 := parse.FieldInfo{
			//Name:   row.ID,
			Type:   row.Title,
			Column: row.Content,
		}
		fmt.Println(str1)
	}

	fmt.Println("===============================")

	fmt.Println(str)
	fmt.Println("===============================")
	//news, err := r.db.FindAllFrom(models.NewsTable, "id", 1)
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
