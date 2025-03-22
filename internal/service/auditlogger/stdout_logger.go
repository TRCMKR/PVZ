package auditlogger

import (
	"fmt"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

func (s *Service) stdoutLogger(batches <-chan models.Log) {
	go func() {
		for batch := range batches {
			fmt.Print(batch.String())
		}
	}()
}
