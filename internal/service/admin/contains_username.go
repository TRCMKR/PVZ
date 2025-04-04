package admin

import (
	"context"
)

// ContainsUsername ...
func (s *Service) ContainsUsername(ctx context.Context, username string) (bool, error) {
	return s.Storage.ContainsUsername(ctx, username)
}
