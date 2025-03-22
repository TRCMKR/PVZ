package auditlogger

import (
	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

func (s *Service) filter(inputChannel <-chan models.Log, operator func(log models.Log) bool) <-chan models.Log {
	outputChannel := make(chan models.Log)
	go func() {
		defer close(outputChannel)
		for log := range inputChannel {
			if !operator(log) {
				continue
			}
			outputChannel <- log
		}
	}()

	return outputChannel
}
