package app

import (
	"log/slog"

	restapp "github.com/Oxeeee/bank-microservices/billing/internal/app/rest"
	"github.com/Oxeeee/bank-microservices/billing/internal/service"
)

type App struct {
	RESTSrv *restapp.App
}

func New(log *slog.Logger, restPort int, service service.BillingService) *App {
	restapp := restapp.New(log, restPort, service)

	return &App{
		RESTSrv: restapp,
	}
}
