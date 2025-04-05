package kafka

import (
	"context"
	"log"
	"time"

	"github.com/IBM/sarama"
	"golang.org/x/sync/errgroup"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

type logsStorage interface {
	GetAndMarkLogs(context.Context, int) ([]models.Log, error)
	UpdateLog(context.Context, int, int, int) error
}

// Pool that works with Kafka
type Pool struct {
	storage logsStorage
	jobs    chan models.Log
	done    chan models.Log
	failed  chan models.Log
}

// NewPool create instance of Pool that works with Kafka
func NewPool(ctx context.Context, cfg config.Config, interval time.Duration,
	storage logsStorage, bufferSize int) *Pool {
	jobs := make(chan models.Log, bufferSize)
	done := make(chan models.Log, bufferSize)
	failed := make(chan models.Log, bufferSize)

	g, gCtx := errgroup.WithContext(ctx)
	go func() {
		<-gCtx.Done()
		close(jobs)
		close(failed)
		close(done)
	}()

	g.Go(func() error {
		return dbReader(gCtx, interval, storage, 10, jobs)
	})
	g.Go(func() error {
		return producer(gCtx, cfg, jobs, failed)
	})
	g.Go(func() error {
		return updater(gCtx, storage, done, failed)
	})
	g.Go(func() error {
		return consumer(gCtx, cfg, done)
	})

	go func() {
		if err := g.Wait(); err != nil {
			log.Fatalf("Error occurred during kafka pool execution: %v", err)
		}
	}()

	return &Pool{
		storage: storage,
		jobs:    jobs,
		done:    done,
		failed:  failed,
	}
}

func initConsumer(ctx context.Context, cfg config.Config) (sarama.PartitionConsumer, error) {
	consumer, err := sarama.NewConsumer([]string{"localhost:" + cfg.KafkaPort()}, nil)
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
