package routes

import (
	"github.com/gofiber/fiber/v2"
	"go_news_server/internal/handlers"
	"go_news_server/internal/middleware"
)

func PrivateRoutes(a *fiber.App, handler *handlers.NewsHandlers) {
	route := a.Group("/private/")

	route.Get("/list", middleware.JWTProtected(), handler.GetNewsList)
	route.Post("/edit/:Id", middleware.JWTProtected(), handler.EditNewsHandler)
}
