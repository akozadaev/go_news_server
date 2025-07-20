package handlers

import (
	"strconv"

	"go_news_server/internal/models"
	"go_news_server/internal/services"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// CreateCategory создает новую категорию
// POST /categories
func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req models.CategoryCreateRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	// Валидация
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Category name is required",
		})
	}

	if len(req.Name) > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Category name must be less than 100 characters",
		})
	}

	category, err := h.categoryService.CreateCategory(c.Context(), &req)
	if err != nil {
		if _, ok := err.(*models.ValidationError); ok {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create category",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.CategoryResponse{
		Success:  true,
		Category: category,
	})
}

// GetCategoryByID получает категорию по ID
// GET /categories/:id
func (h *CategoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid category ID",
		})
	}

	category, err := h.categoryService.GetCategoryByID(c.Context(), id)
	if err != nil {
		if _, ok := err.(*models.NotFoundError); ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get category",
		})
	}

	return c.JSON(models.CategoryResponse{
		Success:  true,
		Category: category,
	})
}

// GetAllCategories получает все категории с пагинацией
// GET /categories
func (h *CategoryHandler) GetAllCategories(c *fiber.Ctx) error {
	limit := 10
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	categories, total, err := h.categoryService.GetAllCategories(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get categories",
		})
	}

	return c.JSON(models.CategoriesResponse{
		Success:    true,
		Categories: categories,
		Total:      total,
	})
}

// UpdateCategory обновляет категорию
// PUT /categories/:id
func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid category ID",
		})
	}

	var req models.CategoryUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	// Валидация
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Category name is required",
		})
	}

	if len(req.Name) > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Category name must be less than 100 characters",
		})
	}

	category, err := h.categoryService.UpdateCategory(c.Context(), id, &req)
	if err != nil {
		if _, ok := err.(*models.NotFoundError); ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		if _, ok := err.(*models.ValidationError); ok {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update category",
		})
	}

	return c.JSON(models.CategoryResponse{
		Success:  true,
		Category: category,
	})
}

// DeleteCategory удаляет категорию
// DELETE /categories/:id
func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid category ID",
		})
	}

	err = h.categoryService.DeleteCategory(c.Context(), id)
	if err != nil {
		if _, ok := err.(*models.NotFoundError); ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete category",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Category deleted successfully",
	})
}

// GetCategoriesByNewsID получает категории для конкретной новости
// GET /news/:id/categories
func (h *CategoryHandler) GetCategoriesByNewsID(c *fiber.Ctx) error {
	newsIDStr := c.Params("id")
	newsID, err := strconv.ParseInt(newsIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid news ID",
		})
	}

	categories, err := h.categoryService.GetCategoriesByNewsID(c.Context(), newsID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get categories for news",
		})
	}

	return c.JSON(models.CategoriesResponse{
		Success:    true,
		Categories: categories,
		Total:      int64(len(categories)),
	})
}
