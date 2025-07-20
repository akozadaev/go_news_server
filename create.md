

Для создания REST сервера с использованием Fiber и базы данных PostgreSQL, мы начнем с инициализации проекта и настройки необходимых зависимостей. Мы будем использовать Fiber в качестве веб-сервера и


rm для работы с базой данных. Все настройки будут управляться с помощью Viper, а подключение к базе данных будет использовать connection pool.

### Шаг 1: Инициализация проекта

Создадим новый проект и инициализируем его с помощью `go mod`:

```sh
mkdir news-rest-api
cd news-rest-api
go mod init news-rest-api
```

### Шаг 2: Установка зависимостей

Установим необходимые пакеты:

```sh
go get github.com/gofiber/fiber/v2
go get github.com/go-pg/reform
go get github.com/spf13/viper
go get github.com/sirupsen/logrus
go get github.com/jackc/pgx/v4/pgxpool
```

### Шаг 3: Настройка структуры проекта

Создадим структуру директорий для нашего проекта:

```
news-rest-api/
|-- cmd/
|   |-- main.go
|-- config/
|   |-- config.go
|-- handlers/
|   |-- news.go
|-- models/
|   |-- news.go
|-- repository/
|   |-- news.go
|-- services/
|   |-- news.go
|-- utils/
|   |-- logger.go
|-- Dockerfile
|-- go.mod
|-- go.sum
```

### Шаг 4: Конфигурация проекта

Создадим файл конфигурации `config/config.go` для управления настройками через Viper:

```go
package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    DatabaseURL string
    Port        string
}

func LoadConfig() (*Config, error) {
    viper.AddConfigPath(".")
    viper.SetConfigName("config")
    viper.SetConfigType("env")

    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    config := &Config{
        DatabaseURL: viper.GetString("DATABASE_URL"),
        Port:        viper.GetString("PORT"),
    }

    return config, nil
}
```

### Шаг 5: Модели данных

Создадим модели данных в `models/news.go`:

```go
package models

import (
    "gopkg.in/reform.v1"
)

//go:generate reform

// News represents a row in news table.
//reform:News
type News struct {
    ID      int64  `reform:"Id,pk"`
    Title   string `reform:"Title"`
    Content string `reform:"Content"`
}

//go:generate reform

// NewsCategory represents a row in news_categories table.
//reform:NewsCategories
type NewsCategory struct {
    NewsID     int64 `reform:"NewsId,pk"`
    CategoryID int64 `reform:"CategoryId,pk"`
}
```

### Шаг 6: Репозиторий

Создадим репозиторий для работы с базой данных в `repository/news.go`:

```go
package repository

import (
    "context"
    "errors"
    "news-rest-api/models"
    "gopkg.in/reform.v1"
)

type NewsRepository interface {
    GetNews(ctx context.Context) ([]models.News, error)
    UpdateNews(ctx context.Context, news *models.News) error
}

type newsRepository struct {
    db *reform.DB
}

func NewNewsRepository(db *reform.DB) NewsRepository {
    return &newsRepository{db: db}
}

func (r *newsRepository) GetNews(ctx context.Context) ([]models.News, error) {
    var news []models.News
    err := r.db.FindAllFrom(models.NewsTable, "ORDER BY Id", &news)
    if err != nil {
        return nil, err
    }
    return news, nil
}

func (r *newsRepository) UpdateNews(ctx context.Context, news *models.News) error {
    if news.ID == 0 {
        return errors.New("invalid ID")
    }
    _, err := r.db.Update(news)
    return err
}
```

### Шаг 7: Сервисы

Создадим сервисы для управления логикой приложения в `services/news.go`:

```go
package services

import (
    "context"
    "news-rest-api/models"
    "news-rest-api/repository"
)

type NewsService interface {
    GetNews(ctx context.Context) ([]models.News, error)
    UpdateNews(ctx context.Context, news *models.News) error
}

type newsService struct {
    repo repository.NewsRepository
}

func NewNewsService(repo repository.NewsRepository) NewsService {
    return &newsService{repo: repo}
}

func (s *newsService) GetNews(ctx context.Context) ([]models.News, error) {
    return s.repo.GetNews(ctx)
}

func (s *newsService) UpdateNews(ctx context.Context, news *models.News) error {
    return s.repo.UpdateNews(ctx, news)
}
```

### Шаг 8: Обработчики

Создадим обработчики для маршрутов в `handlers/news.go`:

```go
package handlers

import (
    "github.com/gofiber/fiber/v2"
    "news-rest-api/models"
    "news-rest-api/services"
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
```

### Шаг 9: Логирование

Создадим утилиту для логирования в `utils/logger.go`:

```go
package utils

import (
    "github.com/sirupsen/logrus"
)

var Log = logrus.New()

func InitLogger() {
    Log.SetFormatter(&logrus.JSONFormatter{})
    Log.SetLevel(logrus.InfoLevel)
}
```

### Шаг 10: Основной файл

Создадим основной файл `cmd/main.go`, который будет запускать сервер:

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/jackc/pgx/v4/pgxpool"
    "gopkg.in/reform.v1"
    "gopkg.in/reform.v1/dialects/postgresql"
    "news-rest-api/config"
    "news-rest-api/handlers"
    "news-rest-api/repository"
    "news-rest-api/services"
    "news-rest-api/utils"
    "log"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("could not load config: %v", err)
    }

    utils.InitLogger()

    pool, err := pgxpool.Connect(ctx, cfg.DatabaseURL)
    if err != nil {
        utils.Log.Fatalf("could not connect to database: %v", err)
    }
    defer pool.Close()

    db := reform.NewDB(pool, postgresql.Dialect, reform.NewPrintfLogger(log.Printf))

    newsRepo := repository.NewNewsRepository(db)
    newsService := services.NewNewsService(newsRepo)
    newsHandler := handlers.NewNewsHandler(newsService)

    app := fiber.New()

    app.Get("/list", newsHandler.GetNewsList)
    app.Post("/edit/:Id", newsHandler.EditNews)

    log.Fatal(app.Listen(":" + cfg.Port))
}
```

### Шаг 11: Dockerfile

Создадим Dockerfile для упаковки сервиса в контейнер:

```dockerfile
FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /news-rest-api

EXPOSE 8080

CMD ["/news-rest-api"]
```

Теперь, для сборки и запуска контейнера:

```sh
docker build -t news-rest-api .
docker run -p 8080:8080 --env-file .env news-rest-api
```

### Шаг 12: Настройка окружения

Создадим файл `.env` для хранения переменных окружения:

```
DATABASE_URL=postgres://user:password@localhost:5432/news_db
PORT=8080
```

### Заключение

Мы создали REST сервер с использованием Fiber и базы данных PostgreSQL. Реализовали маршруты для редактирования и получения списка новостей, настроили логирование и управление конфигурациями через Viper. Добавили Dockerfile для упаковки сервиса в контейнер.

Вы можете расширить проект, добавив авторизацию, валидацию полей, пагинацию и более детализированную обработку ошибок, как указано в дополнительных пожеланиях.