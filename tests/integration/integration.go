package integration

import (
	"context"
	"database/sql"
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

// InitPostgresContainer ...
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

func InitKafkaContainer(ctx context.Context, cfg config.Config) (testcontainers.Container, testcontainers.Container) {
	network, err := network.New(ctx)
	if err != nil {
		log.Fatalf("could not create network: %v", err)
	}

	zookeeper, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "confluentinc/cp-zookeeper",
			ExposedPorts: []string{"2181/tcp"},
			Env: map[string]string{
				"ZOOKEEPER_CLIENT_PORT": "2181",
			},
			Networks:       []string{network.Name},
			NetworkAliases: map[string][]string{network.Name: {"zookeeper"}},
			WaitingFor:     wait.ForListeningPort("2181/tcp"),
		},
		Started: true,
	})
	if err != nil {
		log.Fatalf("could not start zookeeper: %v", err)
	}

	kafka, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "confluentinc/cp-kafka",
			ExposedPorts: []string{cfg.KafkaPort() + "/tcp"},
			Env: map[string]string{
				"KAFKA_BROKER_ID":                        "1",
				"KAFKA_ZOOKEEPER_CONNECT":                "zookeeper:2181",
				"KAFKA_ADVERTISED_LISTENERS":             "INTERNAL://kafka:29092,EXTERNAL://localhost:" + cfg.KafkaPort(),
				"KAFKA_LISTENERS":                        "INTERNAL://kafka:29092,EXTERNAL://:" + cfg.KafkaPort(),
				"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP":   "INTERNAL:PLAINTEXT,EXTERNAL:PLAINTEXT",
				"KAFKA_INTER_BROKER_LISTENER_NAME":       "INTERNAL",
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
		log.Fatalf("could not start kafka: %v", err)
	}

	return kafka, zookeeper
}
