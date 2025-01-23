package main

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go_news_server/internal/handlers"
	"go_news_server/internal/middleware"
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
	"time"

	_ "context"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	serverConfig, err := config.Load()
	if err != nil {
		fmt.Println(err)
		log.Error().Stack().Err(err)
	}

	loggerLevel := zapcore.Level(serverConfig.LoggingConfig.Level)
	if !serverConfig.LoggingConfig.Development {
		loggerLevel = zapcore.ErrorLevel
	}

	logging.SetConfig(&logging.Config{
		Encoding:        serverConfig.LoggingConfig.Encoding,
		Level:           loggerLevel,
		InfoFilename:    serverConfig.LoggingConfig.InfoFilename,
		InfoMaxSize:     serverConfig.LoggingConfig.InfoMaxSize,
		InfoMaxBackups:  serverConfig.LoggingConfig.InfoMaxBackups,
		InfoMaxAge:      serverConfig.LoggingConfig.InfoMaxAge,
		InfoCompress:    serverConfig.LoggingConfig.InfoCompress,
		ErrorFilename:   serverConfig.LoggingConfig.ErrorFilename,
		ErrorMaxSize:    serverConfig.LoggingConfig.ErrorMaxSize,
		ErrorMaxBackups: serverConfig.LoggingConfig.ErrorMaxBackups,
		ErrorMaxAge:     serverConfig.LoggingConfig.ErrorMaxAge,
		ErrorCompress:   serverConfig.LoggingConfig.ErrorCompress,
	})
	defer logging.DefaultLogger().Sync()

	app := fx.New(
		fx.Supply(serverConfig),
		fx.Supply(logging.DefaultLogger().Desugar()),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log.Named("fx")}
		}),
		fx.StopTimeout(serverConfig.ServerConfig.GracefulShutdown*time.Second),
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

	if err := godotenv.Load(".env"); err != nil {
		logger.Errorw("Dont reading .env file")
	}

	config := config.FiberConfig()
	app := fiber.New(config)

	// Get *sql.DB as usual. PostgreSQL example:
	sqlDB, err := sql.Open("postgres", cfg.DBConfig.DataSourceName)
	if err != nil {
		logger.Errorw("Dont open DB", err)
	}
	defer sqlDB.Close()

	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(log.Printf))
	if err != nil {
		logger.Errorw("Dont DB connect")
	}
	newsRepo := repository.NewNewsRepository(db)
	newsService := services.NewNewsService(newsRepo)
	newsHandler := handlers.NewNewsHandler(newsService)
	middleware.FiberMiddleware(app)
	//middleware.JWTProtected()

	routes.PublicRoutes(app, newsHandler)
	routes.NotFoundRoute(app)

	// Start server (with or without graceful shutdown).
	if os.Getenv("STAGE_STATUS") == "dev" {
		utils.StartServer(app, logger)
	} else {
		utils.StartServerWithGracefulShutdown(app, logger)
	}
	logger.Infof("Start to rest api server :%d", cfg.ServerConfig.Port)
	return app
}
