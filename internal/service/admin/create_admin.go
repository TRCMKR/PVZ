package admin

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

// CreateAdmin creates admin
func (s *Service) CreateAdmin(ctx context.Context, admin models.Admin) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.CreateAdmin")
	defer span.Finish()

	ok, err := s.ContainsUsername(ctx, admin.Username)
	if err != nil {
		span.SetTag("error", err)

		return err
	}
	if ok {
		s.logger.Error(ErrUsernameUsed.Error(),
			zap.String("username", admin.Username),
			zap.Error(ErrUsernameUsed),
		)
		span.SetTag("error", ErrUsernameUsed)

		return ErrUsernameUsed
	}

	ok, err = s.ContainsID(ctx, admin.ID)
	if err != nil {
		span.SetTag("error", err)

		return err
	}
	if ok {
		s.logger.Error(ErrIDUsed.Error(),
			zap.Int("id", admin.ID),
			zap.Error(ErrIDUsed),
		)
		span.SetTag("error", ErrIDUsed)

		return ErrIDUsed
	}

	return s.Storage.CreateAdmin(ctx, admin)
}
