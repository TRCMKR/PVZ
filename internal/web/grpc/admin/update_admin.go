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

// UpdateAdmin is a grpc handler over service for updating admin
func (h *Handler) UpdateAdmin(ctx context.Context, req *proto.UpdateAdminRequest) (*proto.UpdateAdminResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "handler.UpdateAdmin")
	defer span.Finish()

	logger := h.logger.With(
		zap.String("handler", "UpdateAdmin"),
	)

	logger.Info("Received request to update admin",
		zap.String("username", req.GetUsername()),
	)

	if req.GetUsername() == "" || req.GetPassword() == "" || req.GetNewPassword() == "" {
		logger.Error(errMissingFields.Error(),
			zap.String("username", req.GetUsername()),
			zap.Error(errMissingFields),
		)
		span.SetTag("error", errMissingFields)

		return nil, errMissingFields
	}

	admin := *models.NewAdmin(0, req.GetUsername(), req.GetNewPassword())
	span.SetTag("created username", admin.Username)

	err := h.Service.UpdateAdmin(ctx, req.GetUsername(), req.GetPassword(), admin)
	if err != nil {
		span.SetTag("error", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Info("Successfully updated admin",
		zap.String("username", req.GetUsername()),
	)
	span.SetTag("success username", admin.Username)

	return &proto.UpdateAdminResponse{
		Output: "success",
	}, nil
}
