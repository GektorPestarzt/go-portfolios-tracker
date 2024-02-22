package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go-portfolios-tracker/internal/auth"
	"go-portfolios-tracker/internal/config"
	"go-portfolios-tracker/internal/logging"
	"go-portfolios-tracker/pkg/repository/sqlite"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func main() {
	os.Setenv("CONFIG_PATH", "./config.yaml")
	cfg := config.MustLoad()
	fmt.Println(cfg)

	logging.MustLoad(cfg.Env)
	logger := logging.GetLogger()

	logger.Info("create database")
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logger.Fatal("failed to init repository", err)
	}

	_ = storage

	logger.Info("create router")
	router := httprouter.New()

	logger.Info("register auth handler")
	handler := auth.NewHandler(logger)
	handler.Register(router)

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
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

func mustCreateSocketListener(logger *logging.Logger) net.Listener {
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
