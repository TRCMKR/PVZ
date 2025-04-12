package order

import (
	"context"

	"github.com/opentracing/opentracing-go"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

// Returns gets all returned orders
func (s *Service) Returns(ctx context.Context) ([]models.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.Returns")
	defer span.Finish()

	return s.Storage.GetReturns(ctx, nil)
}
