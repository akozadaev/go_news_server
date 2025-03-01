package repository

import (
	"context"
	"go_news_server/internal/models"
	"gopkg.in/reform.v1"
)

type NewsRepository struct {
	DB *reform.DB
}

func (r *NewsRepository) UpdateNews(ctx context.Context, news *models.News) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if news.Title != "" {
		_, err = tx.ExecContext(ctx, "UPDATE News SET Title = $1 WHERE Id = $2", news.Title, news.Id)
		if err != nil {
			return err
		}
	}

	if news.Content != "" {
		_, err = tx.ExecContext(ctx, "UPDATE News SET Content = $1 WHERE Id = $2", news.Content, news.Id)
		if err != nil {
			return err
		}
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM NewsCategories WHERE news_id = $1", news.Id)
	if err != nil {
		return err
	}

	for _, category := range news.Categories {
		_, err = tx.ExecContext(ctx, "INSERT INTO NewsCategories (news_Id, category_id) VALUES ($1, $1)", news.Id, category)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *NewsRepository) GetNewsList(ctx context.Context, limit, offset int) ([]models.News, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT n.Id, n.Title, n.Content, 
               array_agg(nc.category_id) as Categories
        FROM News n
        LEFT JOIN NewsCategories nc ON n.Id = nc.news_id
        GROUP BY n.Id
        LIMIT $1 OFFSET $2`, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newsList []models.News
	for rows.Next() {
		var news models.News
		var categories string
		if err := rows.Scan(&news.Id, &news.Title, &news.Content, &categories); err != nil {
			return nil, err
		}

		newsList = append(newsList, news)
	}

	return newsList, nil
}

func (r *NewsRepository) AddNews(ctx context.Context, news *models.News) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if news.Title != "" && news.Content != "" {
		_, err = tx.ExecContext(ctx,
			"insert into news (title, content) VALUES ($1, $2)",
			news.Title, news.Content)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
