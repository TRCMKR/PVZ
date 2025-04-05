package kafka

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	"gitlab.ozon.dev/alexplay1224/homework/internal/models"

	"github.com/IBM/sarama"
)

func producer(ctx context.Context, cfg config.Config, jobs <-chan models.Log, failed chan<- models.Log) error {
	producer, err := sarama.NewSyncProducer([]string{"localhost:" + cfg.KafkaPort()}, nil)
	if err != nil {
		return err
	}
	defer producer.Close()

	log.Print("Starting producer")
	for job := range jobs {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		jobBytes, err := json.Marshal(job)
		if err != nil {
			log.Printf("Error marshaling JSON: %v\n", err)

			continue
		}

		message := &sarama.ProducerMessage{
			Topic: "logs",
			Key:   sarama.StringEncoder(strconv.Itoa(job.ID)),
			Value: sarama.ByteEncoder(jobBytes),
		}

		_, _, err = producer.SendMessage(message)
		if err != nil {
			failed <- job
		}
	}

	return nil
}
