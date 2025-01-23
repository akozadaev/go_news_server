package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go_news_server/internal/models"
	"go_news_server/internal/services"
	"strconv"
)

type NewsHandler struct {
	service services.NewsService
}

func NewNewsHandler(service services.NewsService) *NewsHandler {
	return &NewsHandler{service: service}
}

func (h *NewsHandler) GetNewsList(c *fiber.Ctx) error {
	news, err := h.service.GetNews(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"Success": true, "News": news})
}

func (h *NewsHandler) EditNews(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("Id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
	}

	var news models.News
	if err := c.BodyParser(&news); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	news.ID = int64(id)
	if err := h.service.UpdateNews(c.Context(), &news); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"Success": true})
}
