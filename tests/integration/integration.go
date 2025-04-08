package integration

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"

	_ "github.com/lib/pq" // lib/pg ...
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// InitPostgresContainer creates an instance of postgres container with applied migrations
func InitPostgresContainer(ctx context.Context, cfg config.Config) (string, *postgres.PostgresContainer, error) {
	emptyLogger := log.New(io.Discard, "", 0)

	pgContainer, err := postgres.Run(ctx, "postgres:14-alpine",
		postgres.WithDatabase(cfg.DBName()),
		postgres.WithUsername(cfg.Username()),
		postgres.WithPassword(cfg.Password()),
		testcontainers.WithLogger(emptyLogger),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return "", nil, err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return "", nil, err
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return "", nil, err
	}
	defer db.Close()

	time.Sleep(2 * time.Second)

	goose.SetLogger(emptyLogger)
	err = goose.SetDialect("postgres")
	if err != nil {
		return "", nil, err
	}
	rootDir, err := config.GetRootDir()
	if err != nil {
		return "", nil, err
	}
	if err = goose.Up(db, rootDir+"/migrations"); err != nil {
		log.Panicf("failed to run migrations: %v", err)
	}

	return connStr, pgContainer, nil
}

// InitKafkaContainer is used to create testcontainers Kafka with specified config
func InitKafkaContainer(ctx context.Context, cfg config.Config) (testcontainers.Container, error) {
	network, err := network.New(ctx)
	if err != nil {
		log.Fatalf("could not create network: %v", err)
	}

	kafka, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "confluentinc/cp-kafka",
			ExposedPorts: []string{cfg.KafkaPort() + "/tcp"},
			Env: map[string]string{
				"KAFKA_BROKER_ID":                "1",
				"KAFKA_NODE_ID":                  "1",
				"CLUSTER_ID":                     "local_cluster1",
				"KAFKA_KRAFT_MODE":               "true",
				"KAFKA_PROCESS_ROLES":            "controller,broker",
				"KAFKA_CONTROLLER_QUORUM_VOTERS": "1@kafka:9093",
				"KAFKA_ADVERTISED_LISTENERS": fmt.Sprintf("INTERNAL://kafka:29092,EXTERNAL://%s:%s",
					cfg.KafkaHost(), cfg.KafkaPort()),
				"KAFKA_LISTENERS": fmt.Sprintf("INTERNAL://kafka:29092,EXTERNAL://:%s,CONTROLLER://:9093",
					cfg.KafkaPort()),
				"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP":   "INTERNAL:PLAINTEXT,EXTERNAL:PLAINTEXT,CONTROLLER:PLAINTEXT",
				"KAFKA_INTER_BROKER_LISTENER_NAME":       "INTERNAL",
				"KAFKA_CONTROLLER_LISTENER_NAMES":        "CONTROLLER",
				"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR": "1",
			},
			Networks:       []string{network.Name},
			NetworkAliases: map[string][]string{network.Name: {"kafka"}},
			WaitingFor:     wait.ForListeningPort(nat.Port(cfg.KafkaPort() + "/tcp")),
			HostConfigModifier: func(hc *container.HostConfig) {
				hc.PortBindings = nat.PortMap{
					nat.Port(cfg.KafkaPort() + "/tcp"): []nat.PortBinding{
						{
							HostIP:   "0.0.0.0",
							HostPort: cfg.KafkaPort(),
						},
					},
				}
			},
		},
		Started: true,
	})
	if err != nil {
		return nil, err
	}

	return kafka, nil
}
