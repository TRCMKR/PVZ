package admin

import (
	"context"
)

// ContainsID ...
func (s *Service) ContainsID(ctx context.Context, id int) (bool, error) {
	return s.Storage.ContainsID(ctx, id)
}
