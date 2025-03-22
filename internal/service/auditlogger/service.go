package auditlogger

import (
	"context"
	"strings"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

type auditLoggerStorage interface {
	CreateLog(context.Context, models.Log) error
}

type Service struct {
	Storage auditLoggerStorage
	jobs    chan models.Log
}

func NewService(ctx context.Context, logs auditLoggerStorage, workerCount int, batchSize int,
	timeout time.Duration) *Service {
	jobs := make(chan models.Log, batchSize*20)
	go func() {
		defer close(jobs)
		<-ctx.Done()
	}()

	rootDir, err := config.GetRootDir()
	if err != nil {
		panic(err)
	}
	word, err := config.ReadFirstFileWord(rootDir + "/logger.config")
	if err != nil {
		panic(err)
	}
	s := &Service{
		Storage: logs,
		jobs:    jobs,
	}

	dbWorkerCount := workerCount/2 + workerCount%2
	for i := 0; i < dbWorkerCount; i++ {
		go s.dbWorker(ctx, batchSize, timeout, jobs)
	}

	stdoutWorkerCount := workerCount / 2
	operator := func(log models.Log) bool {
		return strings.Contains(log.String(), word)
	}
	for i := 0; i < stdoutWorkerCount; i++ {
		go s.stduoutWorker(batchSize, timeout, jobs, operator)
	}

	return s
}
