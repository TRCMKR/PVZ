package order

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.ozon.dev/alexplay1224/homework/pkg/api/order/proto"
)

// DeleteOrder is grpc handler over service for deleting order
func (h *Handler) DeleteOrder(ctx context.Context, req *proto.DeleteOrderRequest) (*proto.DeleteOrderResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "handler.DeleteOrder")
	defer span.Finish()

	logger := h.logger.With(
		zap.String("handler", "DeleteOrder"),
	)

	logger.Info("Received request to delete order",
		zap.Int("order_id", int(req.GetId())),
	)

	if req.GetId() == 0 {
		logger.Error(errMissingFields.Error(),
			zap.Int("order_id", int(req.GetId())),
			zap.Error(errMissingFields),
		)
		span.SetTag("error", errMissingFields)

		return nil, errMissingFields
	}

	err := h.Service.ReturnOrder(ctx, int(req.GetId()))
	if err != nil {
		span.SetTag("error", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Info("Successfully deleted order",
		zap.Int("order_id", int(req.GetId())),
	)

	return &proto.DeleteOrderResponse{
		Output: "success",
	}, nil
}
