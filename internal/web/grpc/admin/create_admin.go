package admin

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"gitlab.ozon.dev/alexplay1224/homework/pkg/api/admin/proto"
)

// CreateAdmin is a grpc handler over service for creating admin
func (h *Handler) CreateAdmin(ctx context.Context, req *proto.CreateAdminRequest) (*proto.CreateAdminResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "handler.CreateAdmin")
	defer span.Finish()

	logger := h.logger.With(
		zap.String("handler", "CreateAdmin"),
	)

	logger.Info("Received request to create admin",
		zap.String("username", req.GetUsername()),
	)

	if req.GetId() == 0 || req.GetUsername() == "" || req.GetPassword() == "" {
		logger.Error(errMissingFields.Error(),
			zap.Int("admin_id", int(req.GetId())),
			zap.String("username", req.GetUsername()),
			zap.Error(errMissingFields),
		)
		span.SetTag("error", errMissingFields)

		return nil, errMissingFields
	}

	admin := *models.NewAdmin(int(req.GetId()), req.GetUsername(), req.GetPassword())
	span.SetTag("admin_id", int(req.GetId()))

	err := h.Service.CreateAdmin(ctx, admin)
	if err != nil {
		span.SetTag("error", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Info("Successfully created admin",
		zap.String("username", req.GetUsername()),
	)

	return &proto.CreateAdminResponse{
		Output: "success",
	}, nil
}
