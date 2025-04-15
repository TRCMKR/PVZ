package order

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"gitlab.ozon.dev/alexplay1224/homework/internal/query"
	"gitlab.ozon.dev/alexplay1224/homework/pkg/api/order/proto"
)

// GetOrders is grpc handler over service for getting orders
func (h *Handler) GetOrders(ctx context.Context, req *proto.GetOrdersRequest) (*proto.GetOrdersResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "handler.GetOrders")
	defer span.Finish()

	logger := h.logger.With(
		zap.String("handler", "GetOrders"),
	)

	conds := makeConditions(req)
	logger.Info("Received request to get orders",
		zap.Any("conditions", conds),
	)

	orders, err := h.Service.GetOrders(ctx, conds, int(req.GetCount()), int(req.GetPage()))
	if err != nil {
		span.SetTag("error", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	ordersResponse := make([]*proto.Order, 0, len(orders))
	for _, o := range orders {
		ordersResponse = append(ordersResponse, &proto.Order{
			Id:             int32(o.ID),
			UserId:         int32(o.UserID),
			Weight:         o.Weight,
			Price:          o.Price.Amount(),
			Packaging:      int32(o.Packaging),
			ExtraPackaging: int32(o.ExtraPackaging),
			Status:         int32(o.Status),
			ArrivalDate:    timestamppb.New(o.ArrivalDate),
			ExpiryDate:     timestamppb.New(o.ExpiryDate),
			LastChange:     timestamppb.New(o.LastChange),
		})
	}

	logger.Info("Successfully got orders",
		zap.Any("conditions", conds),
	)

	return &proto.GetOrdersResponse{
		Orders: ordersResponse,
	}, nil
}

//nolint:gocognit
//nolint:gocyclo
func makeConditions(req *proto.GetOrdersRequest) []query.Cond {
	conds := make([]query.Cond, 0, 17)

	if req.Id != nil {
		conds = append(conds, query.Cond{
			Field:    "id",
			Value:    req.GetId(),
			Operator: query.Equals,
		})
	}

	if req.UserId != nil {
		conds = append(conds, query.Cond{
			Field:    "user_id",
			Value:    req.GetUserId(),
			Operator: query.Equals,
		})
	}

	if req.Weight != nil {
		conds = append(conds, query.Cond{
			Field:    "weight",
			Value:    req.GetWeight(),
			Operator: query.Equals,
		})
	}

	if req.WeightTo != nil {
		conds = append(conds, query.Cond{
			Field:    "weight",
			Value:    req.GetWeightTo(),
			Operator: query.LessEqualThan,
		})
	}

	if req.WeightFrom != nil {
		conds = append(conds, query.Cond{
			Field:    "weight",
			Value:    req.GetWeightFrom(),
			Operator: query.GreaterEqualThan,
		})
	}

	if req.Price != nil {
		conds = append(conds, query.Cond{
			Field:    "price",
			Value:    req.GetPrice(),
			Operator: query.Equals,
		})
	}

	if req.PriceTo != nil {
		conds = append(conds, query.Cond{
			Field:    "price",
			Value:    req.GetPriceTo(),
			Operator: query.LessEqualThan,
		})
	}

	if req.PriceFrom != nil {
		conds = append(conds, query.Cond{
			Field:    "price",
			Value:    req.GetPriceFrom(),
			Operator: query.GreaterEqualThan,
		})
	}

	if req.Status != nil {
		conds = append(conds, query.Cond{
			Field:    "status",
			Value:    req.GetStatus(),
			Operator: query.Equals,
		})
	}

	if req.ArrivalDate != nil {
		conds = append(conds, query.Cond{
			Field:    "arrival_date",
			Value:    req.GetArrivalDate(),
			Operator: query.Equals,
		})
	}

	if req.ArrivalDateTo != nil {
		conds = append(conds, query.Cond{
			Field:    "arrival_date",
			Value:    req.GetArrivalDateTo(),
			Operator: query.LessEqualThan,
		})
	}

	if req.ArrivalDateFrom != nil {
		conds = append(conds, query.Cond{
			Field:    "arrival_date",
			Value:    req.GetArrivalDateFrom(),
			Operator: query.GreaterEqualThan,
		})
	}

	if req.ExpiryDate != nil {
		conds = append(conds, query.Cond{
			Field:    "expiry_date",
			Value:    req.GetExpiryDate(),
			Operator: query.Equals,
		})
	}

	if req.ExpiryDateTo != nil {
		conds = append(conds, query.Cond{
			Field:    "expiry_date",
			Value:    req.GetExpiryDateTo(),
			Operator: query.LessEqualThan,
		})
	}

	if req.ExpiryDateFrom != nil {
		conds = append(conds, query.Cond{
			Field:    "expiry_date",
			Value:    req.GetExpiryDateFrom(),
			Operator: query.GreaterEqualThan,
		})
	}

	return conds
}
