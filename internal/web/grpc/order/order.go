package order

import (
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.ozon.dev/alexplay1224/homework/internal/service/order"
	"gitlab.ozon.dev/alexplay1224/homework/pkg/api/order/proto"
)

// Handler is a gRPC order handler implementation
type Handler struct {
	Service order.Service
	proto.UnimplementedOrderServiceServer
	logger *zap.Logger
}

var (
	errMissingFields   = status.Errorf(codes.InvalidArgument, "missing required fields")
	errNoSuchPackaging = status.Errorf(codes.InvalidArgument, "no such packaging")
)

// NewHandler creates an instance of new grpc order Handler
func NewHandler(logger *zap.Logger, service order.Service) *Handler {
	return &Handler{
		Service: service,
		logger:  logger,
	}
}
