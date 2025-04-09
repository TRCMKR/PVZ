package admin

import (
	"context"
)

// ContainsID checks if admin is present
func (s *Service) ContainsID(ctx context.Context, id int) (bool, error) {
	return s.Storage.ContainsID(ctx, id)
}
