package services

import (
	"context"
	"database/sql"
	"go_news_server/internal/models"
	"go_news_server/internal/repository"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

// CreateCategory создает новую категорию
func (s *CategoryService) CreateCategory(ctx context.Context, req *models.CategoryCreateRequest) (*models.Category, error) {
	// Проверяем, существует ли категория с таким именем
	existingCategory, err := s.categoryRepo.GetCategoryByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	if existingCategory != nil {
		return nil, &models.ValidationError{Message: "Category with this name already exists"}
	}

	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	err = s.categoryRepo.CreateCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// GetCategoryByID получает категорию по ID
func (s *CategoryService) GetCategoryByID(ctx context.Context, id int64) (*models.Category, error) {
	category, err := s.categoryRepo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, &models.NotFoundError{Message: "Category not found"}
	}

	return category, nil
}

// GetAllCategories получает все категории с пагинацией
func (s *CategoryService) GetAllCategories(ctx context.Context, limit, offset int) ([]models.Category, int64, error) {
	return s.categoryRepo.GetAllCategories(ctx, limit, offset)
}

// UpdateCategory обновляет категорию
func (s *CategoryService) UpdateCategory(ctx context.Context, id int64, req *models.CategoryUpdateRequest) (*models.Category, error) {
	// Проверяем, существует ли категория
	existingCategory, err := s.categoryRepo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existingCategory == nil {
		return nil, &models.NotFoundError{Message: "Category not found"}
	}

	// Проверяем, не существует ли другая категория с таким именем
	if req.Name != existingCategory.Name {
		categoryWithSameName, err := s.categoryRepo.GetCategoryByName(ctx, req.Name)
		if err != nil {
			return nil, err
		}

		if categoryWithSameName != nil {
			return nil, &models.ValidationError{Message: "Category with this name already exists"}
		}
	}

	// Обновляем категорию
	existingCategory.Name = req.Name
	existingCategory.Description = req.Description

	err = s.categoryRepo.UpdateCategory(ctx, existingCategory)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &models.NotFoundError{Message: "Category not found"}
		}
		return nil, err
	}

	return existingCategory, nil
}

// DeleteCategory удаляет категорию
func (s *CategoryService) DeleteCategory(ctx context.Context, id int64) error {
	// Проверяем, существует ли категория
	existingCategory, err := s.categoryRepo.GetCategoryByID(ctx, id)
	if err != nil {
		return err
	}

	if existingCategory == nil {
		return &models.NotFoundError{Message: "Category not found"}
	}

	err = s.categoryRepo.DeleteCategory(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &models.NotFoundError{Message: "Category not found"}
		}
		return err
	}

	return nil
}

// GetCategoriesByNewsID получает категории для конкретной новости
func (s *CategoryService) GetCategoriesByNewsID(ctx context.Context, newsID int64) ([]models.Category, error) {
	return s.categoryRepo.GetCategoriesByNewsID(ctx, newsID)
}
