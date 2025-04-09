//go:build integration

package kafka

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"gitlab.ozon.dev/alexplay1224/homework/internal/service/auditlogger"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/repository"
	"gitlab.ozon.dev/alexplay1224/homework/tests/integration"
)

func TestKafka_Jobs(t *testing.T) {
	t.Parallel()
	defaultTime := time.Now()

	tests := []struct {
		name         string
		logsToAdd    []models.Log
		expectedLogs []models.Log
	}{
		{
			name:         "no logs to add",
			logsToAdd:    []models.Log{},
			expectedLogs: []models.Log{},
		},
		{
			name: "add NO_ATTEMPTS_LEFT logs",
			logsToAdd: []models.Log{
				{
					AdminID:      1,
					JobStatus:    models.NoAttemptsLeftStatus,
					AttemptsLeft: 0,
				},
				{
					AdminID:      1,
					JobStatus:    models.NoAttemptsLeftStatus,
					AttemptsLeft: 0,
				},
				{
					AdminID:      1,
					JobStatus:    models.NoAttemptsLeftStatus,
					AttemptsLeft: 0,
				},
			},
			expectedLogs: []models.Log{
				{
					ID:           1,
					JobStatus:    models.NoAttemptsLeftStatus,
					AttemptsLeft: 0,
				},
				{
					ID:           2,
					JobStatus:    models.NoAttemptsLeftStatus,
					AttemptsLeft: 0,
				},
				{
					ID:           3,
					JobStatus:    models.NoAttemptsLeftStatus,
					AttemptsLeft: 0,
				},
			},
		},
		{
			name: "add FAILED logs",
			logsToAdd: []models.Log{
				{
					AdminID:      1,
					JobStatus:    models.FailedStatus,
					AttemptsLeft: 1,
				},
				{
					AdminID:      1,
					JobStatus:    models.FailedStatus,
					AttemptsLeft: 1,
				},
				{
					AdminID:      1,
					JobStatus:    models.FailedStatus,
					AttemptsLeft: 1,
				},
			},
			expectedLogs: []models.Log{
				{
					ID:           4,
					JobStatus:    models.DoneStatus,
					AttemptsLeft: 1,
				},
				{
					ID:           5,
					JobStatus:    models.DoneStatus,
					AttemptsLeft: 1,
				},
				{
					ID:           6,
					JobStatus:    models.DoneStatus,
					AttemptsLeft: 1,
				},
			},
		},
		{
			name: "add DONE logs",
			logsToAdd: []models.Log{
				{
					AdminID:      1,
					JobStatus:    models.DoneStatus,
					AttemptsLeft: 3,
				},
				{
					AdminID:      1,
					JobStatus:    models.DoneStatus,
					AttemptsLeft: 3,
				},
				{
					AdminID:      1,
					JobStatus:    models.DoneStatus,
					AttemptsLeft: 3,
				},
			},
			expectedLogs: []models.Log{
				{
					ID:           7,
					JobStatus:    models.DoneStatus,
					AttemptsLeft: 3,
				},
				{
					ID:           8,
					JobStatus:    models.DoneStatus,
					AttemptsLeft: 3,
				},
				{
					ID:           9,
					JobStatus:    models.DoneStatus,
					AttemptsLeft: 3,
				},
			},
		},
		{
			name: "add OLD PROCESSING logs",
			logsToAdd: []models.Log{
				{
					AdminID:      1,
					JobStatus:    models.ProcessingStatus,
					AttemptsLeft: 3,
					UpdatedAt:    defaultTime.Add(-1 * time.Hour),
				},
				{
					AdminID:      1,
					JobStatus:    models.ProcessingStatus,
					AttemptsLeft: 3,
					UpdatedAt:    defaultTime.Add(-1 * time.Hour),
				},
				{
					AdminID:      1,
					JobStatus:    models.ProcessingStatus,
					AttemptsLeft: 3,
					UpdatedAt:    defaultTime.Add(-1 * time.Hour),
				},
			},
			expectedLogs: []models.Log{
				{
					ID:           10,
					JobStatus:    models.DoneStatus,
					AttemptsLeft: 3,
				},
				{
					ID:           11,
					JobStatus:    models.DoneStatus,
					AttemptsLeft: 3,
				},
				{
					ID:           12,
					JobStatus:    models.DoneStatus,
					AttemptsLeft: 3,
				},
			},
		},
		{
			name: "add NEW PROCESSING logs",
			logsToAdd: []models.Log{
				{
					AdminID:      1,
					JobStatus:    models.ProcessingStatus,
					AttemptsLeft: 3,
					UpdatedAt:    time.Now(),
				},
				{
					AdminID:      1,
					JobStatus:    models.ProcessingStatus,
					AttemptsLeft: 3,
					UpdatedAt:    time.Now(),
				},
				{
					AdminID:      1,
					JobStatus:    models.ProcessingStatus,
					AttemptsLeft: 3,
					UpdatedAt:    time.Now(),
				},
			},
			expectedLogs: []models.Log{
				{
					ID:           13,
					JobStatus:    models.DoneStatus,
					AttemptsLeft: 3,
				},
				{
					ID:           14,
					JobStatus:    models.DoneStatus,
					AttemptsLeft: 3,
				},
				{
					ID:           15,
					JobStatus:    models.DoneStatus,
					AttemptsLeft: 3,
				},
			},
		},
		{
			name: "add CREATED logs",
			logsToAdd: []models.Log{
				{
					AdminID:      1,
					JobStatus:    models.CreatedStatus,
					AttemptsLeft: 3,
				},
				{
					AdminID:      1,
					JobStatus:    models.CreatedStatus,
					AttemptsLeft: 3,
				},
				{
					AdminID:      1,
					JobStatus:    models.CreatedStatus,
					AttemptsLeft: 3,
				},
			},
			expectedLogs: []models.Log{
				{
					ID:           16,
					JobStatus:    models.DoneStatus,
					AttemptsLeft: 3,
				},
				{
					ID:           17,
					JobStatus:    models.DoneStatus,
					AttemptsLeft: 3,
				},
				{
					ID:           18,
					JobStatus:    models.DoneStatus,
					AttemptsLeft: 3,
				},
			},
		},
	}

	ctx := context.Background()
	rootDir, err := config.GetRootDir()
	require.NoError(t, err)
	err = config.InitEnv(rootDir + "/.env.test")
	require.NoError(t, err)

	cfg := config.NewConfig()

	connStr, _, err := integration.InitPostgresContainer(t.Context(), cfg)
	require.NoError(t, err)
	db, err := postgres.NewDB(t.Context(), connStr)
	require.NoError(t, err)

	_, err = integration.InitKafkaContainer(t.Context(), cfg)
	require.NoError(t, err)

	logsRepo := repository.NewLogsRepo(db)
	_, err = auditlogger.NewService(ctx, cfg, logsRepo, 1, 1, 1*time.Second)
	require.NoError(t, err)

	currentLogs := make([]models.Log, 0, 18)

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := logsRepo.CreateJob(ctx, tt.logsToAdd)
			require.NoError(t, err)

			currentLogs = append(currentLogs, tt.expectedLogs...)

			time.Sleep(1 * time.Second)

			tmp, err := logsRepo.GetLogs(ctx)
			require.NoError(t, err)

			for j, log := range tmp[i*3:] {
				require.Equal(t, log.ID, currentLogs[j].ID)
				require.Equal(t, log.JobStatus, currentLogs[j].JobStatus)
				require.Equal(t, log.AttemptsLeft, currentLogs[j].AttemptsLeft)
			}
		})
	}
}
