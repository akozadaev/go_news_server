package main

import (
	"context"
	"database/sql"
	"fmt"
	"go_news_server/internal/handlers"
	"go_news_server/internal/repository"
	"go_news_server/internal/routes"
	"go_news_server/internal/server/utils"
	"go_news_server/internal/services"
	"go_news_server/pkg/config"
	"go_news_server/pkg/logging"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
)

var serverCmd = &cobra.Command{
	Use: "server:go_news_server",
	Run: func(cmd *cobra.Command, args []string) {
		runApplication()
	},
}

func main() {
	if err := serverCmd.Execute(); err != nil {
		log.Printf("failed to execute command. err: %v", err)
		os.Exit(1)
	}
}

func runApplication() {
	config, err := config.Load()
	if err != nil {
		fmt.Println(err)
		log.Error().Stack().Err(err)
	}

	loggerLevel := zapcore.Level(config.LogLevel)

	logging.SetConfig(&logging.Config{
		Encoding:        config.LogEncoding,
		Level:           loggerLevel,
		InfoFilename:    config.LogInfoFilename,
		InfoMaxSize:     config.LogInfoMaxSize,
		InfoMaxBackups:  config.LogInfoMaxBackups,
		InfoMaxAge:      config.LogInfoMaxAge,
		InfoCompress:    config.LogInfoCompress,
		ErrorFilename:   config.LogErrorFilename,
		ErrorMaxSize:    config.LogErrorMaxSize,
		ErrorMaxBackups: config.LogErrorMaxBackups,
		ErrorMaxAge:     config.LogErrorMaxAge,
		ErrorCompress:   config.LogErrorCompress,
	})
	defer func() {
		if err := logging.DefaultLogger().Sync(); err != nil {
		}
	}()

	app := fx.New(
		fx.Supply(config),
		fx.Supply(logging.DefaultLogger().Desugar()),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log.Named("fx")}
		}),
		fx.Provide(
			newDatabase,
			newNewsRepository,
			newNewsService,
			newNewsHandler,
			newCategoryRepository,
			newCategoryService,
			newCategoryHandler,
			newServer,
		),
		fx.Invoke(
			setupRoutes,
		),
	)

	app.Run()
}

// newDatabase создает подключение к базе данных
func newDatabase(cfg *config.Config) *reform.DB {
	logger := logging.DefaultLogger()

	sqlDB, err := sql.Open("postgres", cfg.DataSourceName)
	if err != nil {
		logger.Errorw("Failed to open database", err)
		panic(err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.DBConnMaxLifetime) * time.Minute)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		logger.Errorw("Failed to ping database", err)
		panic(err)
	}

	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(log.Printf))
	return db
}

// newNewsRepository создает репозиторий для новостей
func newNewsRepository(db *reform.DB) *repository.NewsRepository {
	return &repository.NewsRepository{DB: db}
}

// newNewsService создает сервис для новостей
func newNewsService(repo *repository.NewsRepository) *services.NewsService {
	return &services.NewsService{Repository: repo}
}

// newNewsHandler создает обработчик для новостей
func newNewsHandler(service *services.NewsService) *handlers.NewsHandlers {
	return &handlers.NewsHandlers{Service: service}
}

// newCategoryRepository создает репозиторий для категорий
func newCategoryRepository(db *reform.DB) *repository.CategoryRepository {
	return &repository.CategoryRepository{DB: db}
}

// newCategoryService создает сервис для категорий
func newCategoryService(repo *repository.CategoryRepository) *services.CategoryService {
	return services.NewCategoryService(repo)
}

// newCategoryHandler создает обработчик для категорий
func newCategoryHandler(service *services.CategoryService) *handlers.CategoryHandler {
	return handlers.NewCategoryHandler(service)
}

// setupRoutes настраивает маршруты приложения
func setupRoutes(
	app *fiber.App,
	newsHandler *handlers.NewsHandlers,
	categoryHandler *handlers.CategoryHandler,
	cfg *config.Config,
) {
	routes.PublicRoutes(app, newsHandler)
	routes.PrivateRoutes(app, newsHandler, cfg)
	routes.SetupCategoryRoutes(app, categoryHandler)
	routes.NotFoundRoute(app)
}

func newServer(lc fx.Lifecycle, cfg *config.Config) *fiber.App {
	logger := logging.DefaultLogger()

	fiberConfig := config.FiberConfig(cfg.ServerReadTimeout)
	app := fiber.New(fiberConfig)

	// Настраиваем lifecycle для graceful shutdown
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			// Start server (with or without graceful shutdown).
			if os.Getenv("STAGE_STATUS") == "dev" {
				utils.StartServer(app, logger, cfg)
			} else {
				utils.StartServerWithGracefulShutdown(app, logger, cfg)
			}
			logger.Infof("Start to rest api server :%d", cfg.ServerPort)
			return nil
		},
		OnStop: func(context.Context) error {
			logger.Info("Shutting down server...")
			return app.Shutdown()
		},
	})

	return app
}
