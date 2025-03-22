package auditlogger

import (
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

func (s *Service) batcher(results <-chan models.Log, batchSize int, timeout time.Duration) <-chan models.Log {
	batch := make(chan models.Log, batchSize-1)
	batches := make(chan models.Log, batchSize*20)
	unload := func() {
		currentBatchSize := len(batch)
		for i := 1; i <= currentBatchSize; i++ {
			tmp := <-batch
			batches <- tmp
		}
	}

	ticker := time.NewTicker(timeout)

	go func() {
		defer func() {
			ticker.Stop()
			unload()
			close(batch)
			close(batches)
		}()

		for {

			select {
			case <-ticker.C:
				unload()
			case res, ok := <-results:
				if !ok {
					return
				}

				select {
				case batch <- res:
					//fmt.Println("Пишем в батч, осталось места:", batchSize-len(batch))
				default:
					unload()
					batches <- res
				}
				ticker.Reset(timeout)
			}
		}
	}()

	return batches
}
