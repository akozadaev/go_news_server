package repository

import (
	"context"
	"database/sql"
	"go_news_server/internal/models"

	"gopkg.in/reform.v1"
)

type CategoryRepository struct {
	DB *reform.DB
}

// CreateCategory создает новую категорию
func (r *CategoryRepository) CreateCategory(ctx context.Context, category *models.Category) error {
	query := `
		INSERT INTO "Categories" ("Name", "Description") 
		VALUES ($1, $2) 
		RETURNING "Id", "CreatedAt", "UpdatedAt"
	`

	return r.DB.QueryRowContext(ctx, query, category.Name, category.Description).
		Scan(&category.Id, &category.CreatedAt, &category.UpdatedAt)
}

// GetCategoryByID получает категорию по ID
func (r *CategoryRepository) GetCategoryByID(ctx context.Context, id int64) (*models.Category, error) {
	query := `SELECT "Id", "Name", "Description", "CreatedAt", "UpdatedAt" FROM "Categories" WHERE "Id" = $1`

	var category models.Category
	err := r.DB.QueryRowContext(ctx, query, id).
		Scan(&category.Id, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}

// GetCategoryByName получает категорию по имени
func (r *CategoryRepository) GetCategoryByName(ctx context.Context, name string) (*models.Category, error) {
	query := `SELECT "Id", "Name", "Description", "CreatedAt", "UpdatedAt" FROM "Categories" WHERE "Name" = $1`

	var category models.Category
	err := r.DB.QueryRowContext(ctx, query, name).
		Scan(&category.Id, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}

// GetAllCategories получает все категории с пагинацией
func (r *CategoryRepository) GetAllCategories(ctx context.Context, limit, offset int) ([]models.Category, int64, error) {
	// Получаем общее количество категорий
	var total int64
	err := r.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM "Categories"`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Получаем категории с пагинацией
	query := `
		SELECT "Id", "Name", "Description", "CreatedAt", "UpdatedAt" 
		FROM "Categories" 
		ORDER BY "Name" 
		LIMIT $1 OFFSET $2
	`

	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.Id, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		categories = append(categories, category)
	}

	return categories, total, nil
}

// UpdateCategory обновляет категорию
func (r *CategoryRepository) UpdateCategory(ctx context.Context, category *models.Category) error {
	query := `
		UPDATE "Categories" 
		SET "Name" = $1, "Description" = $2, "UpdatedAt" = CURRENT_TIMESTAMP 
		WHERE "Id" = $3
	`

	result, err := r.DB.ExecContext(ctx, query, category.Name, category.Description, category.Id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	// Обновляем UpdatedAt в модели
	query = `SELECT "UpdatedAt" FROM "Categories" WHERE "Id" = $1`
	return r.DB.QueryRowContext(ctx, query, category.Id).Scan(&category.UpdatedAt)
}

// DeleteCategory удаляет категорию по ID
func (r *CategoryRepository) DeleteCategory(ctx context.Context, id int64) error {
	query := `DELETE FROM "Categories" WHERE "Id" = $1`

	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetCategoriesByNewsID получает категории для конкретной новости
func (r *CategoryRepository) GetCategoriesByNewsID(ctx context.Context, newsID int64) ([]models.Category, error) {
	query := `
		SELECT c."Id", c."Name", c."Description", c."CreatedAt", c."UpdatedAt"
		FROM "Categories" c
		INNER JOIN "NewsCategories" nc ON c."Id" = nc."CategoryId"
		WHERE nc."NewsId" = $1
		ORDER BY c."Name"
	`

	rows, err := r.DB.QueryContext(ctx, query, newsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.Id, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}
