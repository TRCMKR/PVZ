package kafka

import (
	"context"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

//nolint:gocognit
func updater(ctx context.Context, storage logsStorage, done <-chan models.Log, failed <-chan models.Log) error {
	var receivedJob models.Log
	for {
		select {
		case <-ctx.Done():
			return nil
		case job, ok := <-done:
			if !ok {
				continue
			}

			job.Status = models.DoneStatus
			receivedJob = job
		case job, ok := <-failed:
			if !ok {
				continue
			}

			job.Status = models.NoAttemptsLeftStatus
			if job.AttemptsLeft != 0 {
				job.AttemptsLeft--
				job.Status = models.FailedStatus
			}

			receivedJob = job
		}

		err := storage.UpdateLog(ctx, receivedJob.ID, receivedJob.Status, receivedJob.AttemptsLeft)
		if err != nil {
			return err
		}
	}
}
