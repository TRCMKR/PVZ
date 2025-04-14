package order

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.ozon.dev/alexplay1224/homework/pkg/api/order/proto"
	"gitlab.ozon.dev/alexplay1224/homework/pkg/monitoring"
)

// UpdateOrder is grpc handler over service for updating order
func (h *Handler) UpdateOrder(ctx context.Context, req *proto.UpdateOrderRequest) (*proto.UpdateOrderResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "handler.UpdateOrder")
	defer span.Finish()

	logger := h.logger.With(
		zap.String("handler", "GetOrders"),
	)

	logger.Info("Received request to update order",
		zap.Int("orderId", int(req.GetId())),
	)

	if req.GetId() == 0 || req.GetUserId() == 0 || req.GetAction() == "" {
		logger.Error(errMissingFields.Error(),
			zap.Int("order_id", int(req.GetId())),
			zap.Int("user_id", int(req.GetId())),
			zap.String("action", req.GetAction()),
			zap.Error(errMissingFields),
		)
		span.SetTag("error", errMissingFields)

		return nil, errMissingFields
	}

	err := h.Service.ProcessOrder(ctx, int(req.GetUserId()), int(req.GetId()), req.GetAction())
	if err != nil {
		span.SetTag("error", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Info("Successfully updated order",
		zap.Int("orderId", int(req.GetId())),
	)

	if req.GetAction() == "return" {
		monitoring.SetOrdersReturned()
	}

	return &proto.UpdateOrderResponse{
		Output: "success",
	}, nil
}
