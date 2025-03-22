package auditlogger

import (
	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

func (s *Service) CreateLog(order models.Log) {
	s.jobs <- order
}
