package routes

import (
	"github.com/gofiber/fiber/v2"
	"go_news_server/internal/handlers"
)

func PublicRoutes(a *fiber.App, handler *handlers.NewsHandler) {
	route := a.Group("")
	route.Get("/list", handler.GetNewsList)
	route.Post("/edit/:Id", handler.EditNews)
}
