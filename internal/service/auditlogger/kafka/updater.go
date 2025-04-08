package kafka

import (
	"context"
	"log"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

//nolint:gocognit
func updater(ctx context.Context, storage logsStorage, done <-chan models.Log, failed <-chan models.Log) {
	for {
		var ok bool
		var job models.Log

		select {
		case <-ctx.Done():
			return
		case job, ok = <-done:
			if !ok {
				continue
			}

			job.Status = models.DoneStatus
		case job, ok = <-failed:
			if !ok {
				continue
			}

			job.AttemptsLeft--
			if job.AttemptsLeft == 0 {
				job.Status = models.NoAttemptsLeftStatus
			} else {
				job.Status = models.FailedStatus
			}
		}

		err := storage.UpdateLog(ctx, job.ID, job.Status, job.AttemptsLeft)
		if err != nil {
			log.Printf("error updating jog: %v, jogId: %d", err, job.ID)
		}
	}
}
