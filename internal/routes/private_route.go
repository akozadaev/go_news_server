package routes

import (
	"github.com/gofiber/fiber/v2"
	"go_news_server/internal/handlers"
	"go_news_server/internal/middleware"
	"go_news_server/pkg/config"
)

func PrivateRoutes(a *fiber.App, handler *handlers.NewsHandlers, cfg *config.Config) {
	route := a.Group("/private/")

	route.Get("/list", middleware.KeyProtected(cfg.SecretKey), handler.GetNewsList)
	route.Post("/edit/:Id", middleware.KeyProtected(cfg.SecretKey), handler.EditNewsHandler)
}
