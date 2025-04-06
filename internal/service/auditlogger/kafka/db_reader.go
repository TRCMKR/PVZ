package kafka

import (
	"context"
	"log"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

func dbReader(ctx context.Context, interval time.Duration, storage logsStorage,
	batchSize int, jobs chan<- models.Log) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			logs, err := storage.GetAndMarkLogs(ctx, batchSize)
			if err != nil {
				log.Fatal(err)
			}

			for _, log := range logs {
				jobs <- log
			}

			ticker.Reset(interval)
		}
	}
}
