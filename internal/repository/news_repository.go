package repository

import (
	"context"
	"go_news_server/internal/models"
	"strconv"
	"strings"

	"gopkg.in/reform.v1"
)

type NewsRepository struct {
	DB *reform.DB
}

func (r *NewsRepository) UpdateNews(ctx context.Context, news *models.News, categories []int64) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// Update news fields if they are not empty
	if news.Title != "" {
		_, err = tx.ExecContext(ctx, "UPDATE \"News\" SET \"Title\" = $1 WHERE \"Id\" = $2", news.Title, news.Id)
		if err != nil {
			return err
		}
	}

	if news.Content != "" {
		_, err = tx.ExecContext(ctx, "UPDATE \"News\" SET \"Content\" = $1 WHERE \"Id\" = $2", news.Content, news.Id)
		if err != nil {
			return err
		}
	}

	// Delete existing categories
	_, err = tx.ExecContext(ctx, "DELETE FROM \"NewsCategories\" WHERE \"NewsId\" = $1", news.Id)
	if err != nil {
		return err
	}

	// Insert new categories
	for _, category := range categories {
		_, err = tx.ExecContext(ctx, "INSERT INTO \"NewsCategories\" (\"NewsId\", \"CategoryId\") VALUES ($1, $2)", news.Id, category)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *NewsRepository) GetNewsList(ctx context.Context, limit, offset int) ([]models.News, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT n."Id", n."Title", n."Content", 
               COALESCE(array_agg(nc."CategoryId") FILTER (WHERE nc."CategoryId" IS NOT NULL), '{}') as Categories
        FROM "News" n
        LEFT JOIN "NewsCategories" nc ON n."Id" = nc."NewsId"
        GROUP BY n."Id", n."Title", n."Content"
        ORDER BY n."Id"
        LIMIT $1 OFFSET $2`, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newsList []models.News
	for rows.Next() {
		var news models.News
		var categoriesStr string
		if err := rows.Scan(&news.Id, &news.Title, &news.Content, &categoriesStr); err != nil {
			return nil, err
		}

		// Забираем из PostgreSQL категории как массив {1,2,3}
		var categories []int64
		if categoriesStr != "{}" && categoriesStr != "" {
			// Remove { and } and split by comma
			categoriesStr = strings.Trim(categoriesStr, "{}")
			if categoriesStr != "" {
				parts := strings.Split(categoriesStr, ",")
				for _, part := range parts {
					part = strings.TrimSpace(part)
					if part != "" {
						if categoryID, err := strconv.ParseInt(part, 10, 64); err == nil {
							categories = append(categories, categoryID)
						}
					}
				}
			}
		}
		news.Categories = categories

		newsList = append(newsList, news)
	}

	return newsList, nil
}
