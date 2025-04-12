package admin

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

// ContainsID checks if admin is present
func (s *Service) ContainsID(ctx context.Context, id int) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.ContainsID")
	defer span.Finish()

	return s.Storage.ContainsID(ctx, id)
}
