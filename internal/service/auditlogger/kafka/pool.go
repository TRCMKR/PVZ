package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"golang.org/x/sync/errgroup"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

const (
	batchCount = 10
)

type logsStorage interface {
	GetAndMarkLogs(context.Context, int) ([]models.Log, error)
	UpdateLog(context.Context, int, int, int) error
}

// Start starts pool of Kafka related workers
func Start(ctx context.Context, cfg config.Config, interval time.Duration,
	storage logsStorage, batchSize int) {
	if cfg.IsEmpty() {
		return
	}

	jobs := make(chan models.Log, batchSize*batchCount)
	done := make(chan models.Log, batchSize*batchCount)
	failed := make(chan models.Log, batchSize*batchCount)

	g, gCtx := errgroup.WithContext(ctx)
	go func() {
		<-gCtx.Done()
		close(jobs)
		close(failed)
		close(done)
	}()

	g.Go(func() error {
		return dbReader(gCtx, interval, storage, batchSize, jobs)
	})

	g.Go(func() error {
		return jobSender(gCtx, cfg, jobs, failed)
	})

	go updater(gCtx, storage, done, failed)

	g.Go(func() error {
		return logger(gCtx, cfg, done)
	})

	go func() {
		if err := g.Wait(); err != nil {
			log.Fatalf("Error occurred during kafka pool execution: %v", err)
		}
	}()
}

func initConsumer(ctx context.Context, cfg config.Config) (sarama.PartitionConsumer, error) {
	consumer, err := sarama.NewConsumer([]string{fmt.Sprintf("%s:%s", cfg.KafkaHost(), cfg.KafkaPort())},
		sarama.NewConfig())
	if err != nil {
		return nil, err
	}

	partConsumer, err := consumer.ConsumePartition("logs", 0, sarama.OffsetOldest)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		consumer.Close()
		partConsumer.Close()
	}()

	return partConsumer, nil
}
