package admin

import (
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.ozon.dev/alexplay1224/homework/internal/service/admin"
	"gitlab.ozon.dev/alexplay1224/homework/pkg/api/admin/proto"
)

// Handler is a gRPC admin handler implementation
type Handler struct {
	Service admin.Service
	proto.UnimplementedAdminServiceServer
	logger *zap.Logger
}

var (
	errMissingFields = status.Errorf(codes.InvalidArgument, "missing fields")
)

// NewHandler creates an instance of new grpc admin Handler
func NewHandler(logger *zap.Logger, service admin.Service) *Handler {
	return &Handler{
		Service: service,
		logger:  logger,
	}
}
