package main

import (
	"fmt"
	portfoliohttp "go-portfolios-tracker/internal/account/delivery/http"
	portfoliorepo "go-portfolios-tracker/internal/account/repository/sqlite"
	portfoliousecase "go-portfolios-tracker/internal/account/usecase/tinkoff"
	authhttp "go-portfolios-tracker/internal/auth/delivery/http"
	authrepo "go-portfolios-tracker/internal/auth/repository/sqlite"
	authusecase "go-portfolios-tracker/internal/auth/usecase"
	"go-portfolios-tracker/internal/config"
	"go-portfolios-tracker/internal/logging"
	"go-portfolios-tracker/internal/logging/slog"
	"go-portfolios-tracker/pkg/repository/sqlite"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// @title Portfolio Tracker
// @version 0.0.0
// @description Service for collecting and structuring portfolio data of popular Russian investment instruments

// @host localhost:1234
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	os.Setenv("CONFIG_PATH", "./config.yaml")
	cfg := config.MustLoad()
	fmt.Println(cfg)

	logger := slog.NewLogger(cfg.Env)

	logger.Info("create database")
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logger.Fatal("failed to init repository", err)
	}
	userRepo := authrepo.NewUserRepository(storage)
	portfolioRepo := portfoliorepo.NewPortfolioRepository(storage)

	logger.Info("create router")
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	logger.Info("register auth handler")

	authUseCase := authusecase.NewAuthUseCase(userRepo, cfg.Auth.HashSalt, []byte(cfg.Auth.SigningKey), cfg.Auth.TokenTTL)
	authhttp.RegisterHTTPEndpoints(router, authUseCase, logger)

	authMiddleware := authhttp.NewAuthMiddleware(authUseCase, logger)
	api := router.Group("/api", authMiddleware)

	portfolioUseCase := portfoliousecase.NewPortfolioUseCase(logger, portfolioRepo)
	portfoliohttp.RegisterHTTPEndpoints(api, portfolioUseCase, logger)

	run(router, cfg, logger)
}

func run(router *gin.Engine, cfg *config.Config, logger logging.Logger) {
	logger.Info("start application")

	var listener net.Listener
	var err error

	switch cfg.Listen.Type {
	case "sock":
		listener = mustCreateSocketListener(logger)
		logger.Info("server is listening unix socket")
	case "port":
		logger.Info("listen tcp")
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.Host, cfg.Listen.Port))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info(fmt.Sprintf("server is listening port %s:%s", cfg.Listen.Host, cfg.Listen.Port))
	default:
		logger.Fatal(fmt.Sprintf("config: no listen type: %s", cfg.Listen.Type))
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: cfg.Listen.Timeout,
		ReadTimeout:  cfg.Listen.Timeout,
	}

	logger.Fatal(server.Serve(listener))
}

func mustCreateSocketListener(logger logging.Logger) net.Listener {
	appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("create socket")
	socketPath := path.Join(appDir, "app.sock")
	logger.Debug(fmt.Sprintf("socket path: %s", socketPath))

	logger.Info("listen unix socket")
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		logger.Fatal(err)
	}

	return listener
}
