package routes

import (
	"go_news_server/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupCategoryRoutes настраивает маршруты для категорий
func SetupCategoryRoutes(app *fiber.App, categoryHandler *handlers.CategoryHandler) {
	// Группа маршрутов для категорий
	categories := app.Group("/categories")

	// CRUD операции для категорий
	categories.Post("/", categoryHandler.CreateCategory)      // Создание категории
	categories.Get("/", categoryHandler.GetAllCategories)     // Получение всех категорий
	categories.Get("/:id", categoryHandler.GetCategoryByID)   // Получение категории по ID
	categories.Put("/:id", categoryHandler.UpdateCategory)    // Обновление категории
	categories.Delete("/:id", categoryHandler.DeleteCategory) // Удаление категории

	// Дополнительные маршруты
	news := app.Group("/news")
	news.Get("/:id/categories", categoryHandler.GetCategoriesByNewsID) // Получение категорий для новости
}
