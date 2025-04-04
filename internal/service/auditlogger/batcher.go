package auditlogger

import (
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

//nolint:gocognit
func (s *Service) batcher(results <-chan models.Log, batchSize int, timeout time.Duration) <-chan []models.Log {
	batches := make(chan []models.Log, 20)
	ticker := time.NewTicker(timeout)

	go func() {
		defer close(batches)

		batch := make([]models.Log, 0, batchSize)

		unload := func() {
			if len(batch) > 0 {
				batches <- batch
				batch = make([]models.Log, 0, batchSize)
			}
		}

		for {
			select {
			case <-ticker.C:
				unload()
			case res, ok := <-results:
				if !ok {
					unload()
					ticker.Stop()

					return
				}

				batch = append(batch, res)
				if len(batch) == batchSize {
					unload()
				}
				ticker.Reset(timeout)
			}
		}
	}()

	return batches
}
