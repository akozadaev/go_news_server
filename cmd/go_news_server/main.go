package main

import (
	"context"
	"errors"
	"fmt"
	"go_news_server/pkg/config"
	"go_news_server/pkg/logging"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

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
			// setup database
			//database.NewDatabase,
			// server
			newServer,
			//shortenRepository.NewShortenRepository,
			//shorten.NewHandler,
		),
		fx.Invoke(
		//shorten.RouteV1,
		//func(r *gin.Engine) {},
		),
	)
	app.Run()
}

func newServer(lc fx.Lifecycle, cfg *config.Config) /**gin.Engine*/ {
	//gin.SetMode(gin.DebugMode)
	//r := gin.New()
	//
	//r.Use(middleware.CorsMiddleware())
	//r.Use(middleware.TimeoutMiddleware(cfg.ServerConfig.WriteTimeout))
	//r.Use(middleware.LoggingMiddleware())

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.ServerConfig.Port),
		//Handler:      r,
		ReadTimeout:  cfg.ServerConfig.ReadTimeout,
		WriteTimeout: cfg.ServerConfig.WriteTimeout,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logging.FromContext(ctx).Infof("Start to rest api server :%d", cfg.ServerConfig.Port)
			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logging.DefaultLogger().Errorw("failed to close http server", "err", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logging.FromContext(ctx).Info("Stopped rest api server")
			return srv.Shutdown(ctx)
		},
	})
	//return r
}
