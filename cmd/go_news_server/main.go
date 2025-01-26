package main

import (
	_ "context"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go_news_server/internal/handlers"
	"go_news_server/internal/repository"
	"go_news_server/internal/routes"
	"go_news_server/internal/server/utils"
	"go_news_server/internal/services"
	"go_news_server/pkg/config"
	"go_news_server/pkg/logging"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
	_ "net/http/pprof"
	"os"
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
	defer logging.DefaultLogger().Sync()

	app := fx.New(
		fx.Supply(config),
		fx.Supply(logging.DefaultLogger().Desugar()),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log.Named("fx")}
		}),
		fx.Provide(
			newServer,
		),
		fx.Invoke(
			func(r *fiber.App) {},
		),
	)

	app.Run()
}

func newServer(lc fx.Lifecycle, cfg *config.Config) *fiber.App {
	logger := logging.DefaultLogger()

	fiberConfig := config.FiberConfig(cfg.ServerReadTimeout)
	app := fiber.New(fiberConfig)

	// Get *sql.DB as usual. PostgreSQL example:
	sqlDB, err := sql.Open("postgres", cfg.DataSourceName)
	if err != nil {
		logger.Errorw("Dont open DB", err)
	}
	defer sqlDB.Close()

	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(log.Printf))
	if err != nil {
		logger.Errorw("Dont DB connect")
	}

	newsRepo := &repository.NewsRepository{DB: db}
	newsService := &services.NewsService{Repository: newsRepo}
	newsHandler := &handlers.NewsHandlers{Service: newsService}

	routes.PublicRoutes(app, newsHandler)
	routes.PrivateRoutes(app, newsHandler)
	routes.NotFoundRoute(app)

	// Start server (with or without graceful shutdown).
	if os.Getenv("STAGE_STATUS") == "dev" {
		utils.StartServer(app, logger, cfg)
	} else {
		utils.StartServerWithGracefulShutdown(app, logger, cfg)
	}
	logger.Infof("Start to rest api server :%d", cfg.DBPort)

	return app
}
