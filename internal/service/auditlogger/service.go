package auditlogger

import (
	"context"
	"log"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

type auditLoggerStorage interface {
	CreateLog(context.Context, []models.Log) error
}

type Service struct {
	Storage auditLoggerStorage
	jobs    chan models.Log
}

func NewService(ctx context.Context, logs auditLoggerStorage, workerCount int, batchSize int,
	timeout time.Duration) (*Service, error) {
	jobs := make(chan models.Log, batchSize*20*workerCount)

	g, gCtx := errgroup.WithContext(ctx)

	go func() {
		<-gCtx.Done()
		close(jobs)
	}()

	rootDir, err := config.GetRootDir()
	if err != nil {
		return nil, err
	}

	word, err := config.ReadFirstFileWord(rootDir + "/logger.config")
	if err != nil {
		return nil, err
	}

	s := &Service{
		Storage: logs,
		jobs:    jobs,
	}

	dbWorkerCount := workerCount/2 + workerCount%2
	for i := 0; i < dbWorkerCount; i++ {
		g.Go(func() error {
			return s.dbWorker(gCtx, batchSize, timeout, jobs)
		})
	}

	operator := func(log models.Log) bool {
		return strings.Contains(log.String(), word)
	}
	stdoutWorkerCount := workerCount / 2
	for i := 0; i < stdoutWorkerCount; i++ {
		g.Go(func() error {
			return s.stdoutWorker(gCtx, batchSize, timeout, jobs, operator)
		})
	}

	go func() {
		if err = g.Wait(); err != nil {
			log.Fatalf("Error occurred during service execution: %v", err)
		}
	}()

	return s, nil
}
