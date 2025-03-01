package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	loggerinterceptors "github.com/Oxeeee/bank-microservices/billing/pkg/logger_interceptors"
)

type App struct {
	log    *slog.Logger
	server *http.Server
	port   int
}

func New(log *slog.Logger, grpcAddr string, restPort int) *App {
	mux := http.NewServeMux()
	
	
	
	loggedMux := loggerinterceptors.LoggerMiddleware(log)(mux)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", restPort),
		Handler:      loggedMux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return &App{
		server: server,
		port:   restPort,
		log:    log,
	}
}
