package admin

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

// ContainsUsername checks if username is present
func (s *Service) ContainsUsername(ctx context.Context, username string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.ContainsUsername")
	defer span.Finish()

	return s.Storage.ContainsUsername(ctx, username)
}
