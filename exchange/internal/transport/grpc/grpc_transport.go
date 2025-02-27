package grpctransport

import (
	"context"
	"log/slog"

	"github.com/Oxeeee/bank-microservices/exchange/internal/service"
	customerrors "github.com/Oxeeee/bank-microservices/exchange/pkg/custom_errors"
	exchangev1 "github.com/Oxeeee/bank-microservices/exchange/proto/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	exchangev1.UnimplementedExchangeServer
	service service.ExchangeService
	log     *slog.Logger
}

func Register(gRPC *grpc.Server, service service.ExchangeService, log *slog.Logger) {
	exchangev1.RegisterExchangeServer(gRPC, &serverAPI{
		service: service,
		log:     log,
	})
}

func (s *serverAPI) Convert(ctx context.Context, req *exchangev1.ConvertRequest) (*exchangev1.ConvertResponse, error) {
	converted, err := s.service.ExchangeRequest(req.OrigialCurrencyType, req.ConvertedCurrencyType, req.OriginalCurrencyValue)
	if err != nil && err != customerrors.ErrNotFound {
		return nil, status.Error(codes.Internal, "Can not convert value")
	} else if err == customerrors.ErrNotFound {
		return nil, status.Error(codes.NotFound, "Can not found currency type")
	}

	return &exchangev1.ConvertResponse{
		ConvertedCurrencyType:  req.ConvertedCurrencyType,
		ConvertedCurrencyValue: converted,
	}, nil
}
