package handlers

import (
	"context"
	"go_news_server/internal/models"
	"go_news_server/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type NewsHandlers struct {
	Service services.NewsServiceInterface
}

func (h *NewsHandlers) EditNewsHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("Id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID")
	}

	var payload struct {
		Id         int64   `json:"Id"`
		Title      *string `json:"Title"`
		Content    *string `json:"Content"`
		Categories []int64 `json:"Categories"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	news := &models.News{Id: id}
	if payload.Title != nil {
		news.Title = *payload.Title
	}
	if payload.Content != nil {
		news.Content = *payload.Content
	}

	if err := h.Service.UpdateNews(context.Background(), news, payload.Categories); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"Success": true})
}

// GetNewsList возвращает список новостей.
func (h *NewsHandlers) GetNewsList(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	newsList, err := h.Service.GetNewsList(context.Background(), limit, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"Success": true,
		"News":    newsList,
	})
}
