package admin

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.ozon.dev/alexplay1224/homework/pkg/api/admin/proto"
)

// DeleteAdmin is a grpc handler over service for deleting admin
func (h *Handler) DeleteAdmin(ctx context.Context, req *proto.DeleteAdminRequest) (*proto.DeleteAdminResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "handler.DeleteAdmin")
	defer span.Finish()

	logger := h.logger.With(
		zap.String("handler", "DeleteAdmin"),
	)

	logger.Info("Received request to delete admin",
		zap.String("username", req.GetUsername()),
	)

	if req.GetUsername() == "" || req.GetPassword() == "" {
		logger.Error(errMissingFields.Error(),
			zap.String("username", req.GetUsername()),
			zap.Error(errMissingFields),
		)
		span.SetTag("error", errMissingFields)

		return nil, errMissingFields
	}

	err := h.Service.DeleteAdmin(ctx, req.GetPassword(), req.GetUsername())
	if err != nil {
		span.SetTag("error", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Info("Successfully deleted admin",
		zap.String("username", req.GetUsername()),
	)
	span.SetTag("username", req.GetUsername())

	return &proto.DeleteAdminResponse{
		Output: "success",
	}, nil
}
