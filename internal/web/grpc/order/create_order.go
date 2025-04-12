package order

import (
	"context"

	"github.com/Rhymond/go-money"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"gitlab.ozon.dev/alexplay1224/homework/pkg/api/order/proto"
	"gitlab.ozon.dev/alexplay1224/homework/pkg/monitoring"
)

// CreateOrder is grpc handler over service for creating order
func (h *Handler) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "handler.CreateOrder")
	defer span.Finish()

	logger := h.logger.With(
		zap.String("handler", "CreateOrder"),
	)

	logger.Info("Received request to create order",
		zap.Int("order_id", int(req.GetId())),
	)

	if req.GetId() == 0 || req.GetUserId() == 0 || req.GetWeight() == 0 || req.GetPrice() == 0 ||
		req.GetExpiryDate().AsTime().IsZero() {
		logger.Info("Received request to create order",
			zap.Int("order_id", int(req.GetId())),
		)
		span.SetTag("error", errMissingFields)

		return nil, errMissingFields
	}

	packagingName := models.GetPackagingName(models.PackagingType(req.GetPackaging()))
	extraPackagingName := models.GetPackagingName(models.PackagingType(req.GetExtraPackaging()))

	packaging, err := getPackaging(packagingName)
	extraPackaging, errExtra := getPackaging(extraPackagingName)
	if err != nil || errExtra != nil {
		logger.Error(errNoSuchPackaging.Error(),
			zap.Int("order_id", int(req.GetId())),
			zap.Error(errNoSuchPackaging),
		)
		span.SetTag("error", errNoSuchPackaging)

		return nil, errNoSuchPackaging
	}

	packagings := make([]models.Packaging, 0, 2)
	packagings = append(packagings, packaging)
	packagings = append(packagings, extraPackaging)

	err = h.Service.AcceptOrder(ctx, int(req.GetId()), int(req.GetUserId()), req.GetWeight(),
		*money.New(req.GetPrice(), money.RUB), req.GetExpiryDate().AsTime(), packagings)
	if err != nil {
		span.SetTag("error", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Info("Successfully created order",
		zap.Int("order_id", int(req.GetId())),
	)

	monitoring.SetOrdersCreated()
	monitoring.SetPackagingUsage(packagingName)
	monitoring.SetPackagingUsage(extraPackagingName)
	monitoring.SetOrderTotalPrice(req.GetPrice())

	return &proto.CreateOrderResponse{
		Output: "success",
	}, nil
}

func getPackaging(packagingStr string) (models.Packaging, error) {
	var packaging models.Packaging
	if packagingStr != "" {
		packaging = models.GetPackaging(packagingStr)
		if packaging == nil {
			return nil, errNoSuchPackaging
		}
	} else {
		return nil, errNoSuchPackaging
	}

	return packaging, nil
}
