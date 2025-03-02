package restapp

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Oxeeee/bank-microservices/billing/internal/service"
	"github.com/Oxeeee/bank-microservices/billing/internal/transport/api"
	loggerinterceptors "github.com/Oxeeee/bank-microservices/billing/pkg/logger_interceptors"
)

type App struct {
	log    *slog.Logger
	server *http.Server
	port   int
}

func New(log *slog.Logger, restPort int, service service.BillingService) *App {
	billingHandler := api.NewBillingHandler(log, service)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", billingHandler.Register)

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

func (a *App) MustRun() {
	const op = "rest.run"
	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	if err := a.server.ListenAndServe(); err != nil {
		panic(err)
	}

	log.Info("REST server is running")
}
